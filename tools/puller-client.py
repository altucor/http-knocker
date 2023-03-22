import time
import json
import requests


class HttpKnockerPullerClient:
    def __init__(self, host: str, port: int, base_prefix: str):
        self.__base = "http://" + host + ":" + str(port) + base_prefix

    def get_last_updates(self):
        r = requests.get(self.__base + "/getLastUpdates")
        print(f"Status code: {r.status_code}")
        decoded = json.loads(r.text)
        print(f"decoded: {decoded}")
        return decoded

    def accept_updates(self, accepted_rule_ids):
        print(accepted_rule_ids)
        r = requests.post(
            self.__base + "/acceptUpdates", 
            data={
                'accepted_rules': json.dumps(accepted_rule_ids)
        })
        r.text

class IpTablesExecutor:
    def __init__(self, rule):
        self.__rule = rule
        pass

    def execute(self):
        print(f"Executed rule: {self.__rule}")
        return True

def main():
    httpKnocker = HttpKnockerPullerClient("127.0.0.1", 8001, "/puller/test")

    while True:
        time.sleep(1)
        accepted_rules = []
        rules = httpKnocker.get_last_updates()
        if rules is None:
            continue
        print(f"Rules arr size: {len(rules)}")
        for rule in rules:
            # If rule successfully executed then add it's id to accepted list
            if IpTablesExecutor(rule).execute():
                accepted_rules.append(rule["id"])
        httpKnocker.accept_updates(accepted_rules)



if __name__ == "__main__":
    main()
