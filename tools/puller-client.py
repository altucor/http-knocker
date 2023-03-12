import time
import json
import urllib
import http.client


class HttpKnockerPullerClient:
    def __init__(self, host: str, port: int, base_prefix: str):
        self.__host = host
        self.__port = port
        self.__base_prefix = base_prefix

    def do_request(self, scheme: str, request: str, body, headers):
        connection = http.client.HTTPConnection(self.__host + ":" + str(self.__port))
        connection.request(scheme, self.__base_prefix + request, body, headers)
        response = connection.getresponse()
        print("Status: {} and reason: {}".format(response.status, response.reason))
        body = response.read()
        connection.close()
        return body

    def get_last_updates(self):
        response = self.do_request("GET", "/getLastUpdates", "", {})
        decoded = json.loads(response)
        return decoded

    def accept_updates(self, accepted_rule_ids):
        headers = {'Content-type': 'application/json'}
        body = json.dumps({'accepted_rules': accepted_rule_ids})
        print(body)
        response = self.do_request("POST", "/acceptUpdates", body, headers)

def main():
    httpKnocker = HttpKnockerPullerClient("127.0.0.1", 8001, "/puller/test")

    while True:
        time.sleep(1)
        accepted_rules = []
        rules = httpKnocker.get_last_updates()
        # print(json.loads(rules[0]))
        for rule in rules:
            print(rule)
            accepted_rules.append(rule["id"])
        httpKnocker.accept_updates(accepted_rules)



if __name__ == "__main__":
    main()
