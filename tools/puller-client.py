import time
import http.client


class HttpKnockerPullerClient:
    def __init__(self, host: str, port: int, base_prefix: str):
        self.__host = host
        self.__port = port
        self.__base_prefix = base_prefix

    def do_request(self, request: str):
        connection = http.client.HTTPSConnection(self.__host + ":" + str(self.__port))
        connection.request("GET", self.__base_prefix + request)
        response = connection.getresponse()
        print("Status: {} and reason: {}".format(response.status, response.reason))
        connection.close()
        return response

    def pull_updates(self):
        response = self.do_request("updates")
        print(response)


def main():
    httpKnocker = HttpKnockerPullerClient("127.0.0.1", 8001, "puller/test")

    while True:
        time.sleep(1)
        httpKnocker.pull_updates()


if __name__ == "__main__":
    main()
