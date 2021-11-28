# web-framework-tests

A small project testing several Go and Python frameworks

## About

This projects implements the same API in several different backend frameworks, and it aims to test the performance of each framework with minimal configuration.

Currently supported frameworks:

- Vanilla Go (built-in http package)
- Gin (Go)

To-do:

- Echo (Go)
- FastAPI (Python)
- Flask (Python)

Performance is tested for both local data handling and for saving data to a Postgres database.

This project also includes a profiler client (in Python) that tests each framework and compares the results. The project uses simulated battery cell test data (what I'm most familar with :D).

## Setup
- Install Go (≥ 1.17) and Python (≥ 3.8)
- From the `/profiler` directory, run `pip install -r requirements.txt` to install Python dependencies.
- From each framework directory, run `go mod tidy` to install Go dependencies
- For each framework, run `go run cmd/main.go` to start the web server, then from the `/profiler` directory run `python run_test.py` to run the test

## Notes

This should not be considered a standardized benchmark test, as it was mainly intended as a project that measures performance with little configuration. Benchmark tests typically include a number of hardware and software optimizations that result in performance numbers much higher than those achieved here. 

[Tech Empower's Web Framework Benchmarks](https://www.techempower.com/benchmarks/) are a much better source for performance comparisons.