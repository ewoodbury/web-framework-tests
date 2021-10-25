import time
import requests


n_requests = 1000


def main():
    t0 = time.time()
    for i in range(n_requests):
        requests.post("http://127.0.0.1:8000/data",
                      data={"voltage": 3.4 + i*.001, "current": 1.1 * i*.001})
    t1 = time.time()
    print(f"In-memory array: {1000*(t1 - t0)/n_requests:1.2f} ms/request")


if __name__ == "__main__":
    main()
