use std::sync::{Arc, Mutex};
use std::thread;
use std::time::Duration;

use hyper::{Body, Request, Response, Server};
use hyper::service::{make_service_fn, service_fn};

struct Worker {
    num: u32,
}

async fn cpu_intensive_algorithm(n: i32) {
    for i in 1..=n {
        for j in 1..=n {
            for k in 1..=n {
                for l in 1..=n {
                    // Perform some CPU-intensive computation
                    let _ = i * j * k * l;
                }
            }
        }
    }
}

async fn multi_threaded_handler(worker: Arc<Mutex<Worker>>, req: Request<Body>) -> Result<Response<Body>, hyper::Error> {
    // Acquire a worker (lock) from the mutex
    let worker = worker.lock().unwrap();
    let worker_num = worker.num;
    drop(worker); // Release the lock immediately

    // Perform CPU-intensive operations
    cpu_intensive_algorithm(10).await;

    // Simulate work by sleeping for a short duration
    thread::sleep(Duration::from_millis(100));

    // Prepare the response
    let response = format!("request processed by worker {}", worker_num);
    let body = Body::from(response);

    Ok(Response::new(body))
}

#[tokio::main]
async fn main() {
    // Create a worker
    let worker = Arc::new(Mutex::new(Worker { num: 1 }));

    // Create a server
    let make_svc = make_service_fn(|_conn| {
        let worker = worker.clone();
        async {
            Ok::<_, hyper::Error>(service_fn(move |req| {
                multi_threaded_handler(worker.clone(), req)
            }))
        }
    });

    let addr = ([127, 0, 0, 1], 8000).into();
    let server = Server::bind(&addr).serve(make_svc);

    println!("Server listening on http://{}", addr);

    if let Err(e) = server.await {
        eprintln!("Server error: {}", e);
    }
}
