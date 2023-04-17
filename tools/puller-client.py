import time
import json
import requests
import subprocess

class HttpKnockerPullerClient:
    def __init__(self, host: str, port: int, base_prefix: str):
        self.__base = "http://" + host + ":" + str(port) + base_prefix

    def get_last_updates(self):
        print(f"URL: {self.__base}")
        r = requests.get(self.__base + "/getLastUpdates")
        print(f"Status code: {r.status_code}")
        decoded = json.loads(r.text)
        print(f"decoded: {decoded}")
        return decoded

    def accept_updates(self, accepted_rule_ids):
        print(f"accepted ids: {accepted_rule_ids}")
        r = requests.post(
            self.__base + "/acceptUpdates", 
            data={
                'accepted_rules': json.dumps(accepted_rule_ids)
        })
        print(f"accepted response: {r.text}")

    def push_frw_rules(self, rules_set):
        print(f"rule set: {rules_set}")
        r = requests.post(
            self.__base + "/pushRulesSet", 
            data={
                'rules': json.dumps(rules_set)
        })
        print(f"rule set reponse: {r.text}")

def run_shell_cmd(*args):
    result = subprocess.run(
        args,
        capture_output=True
    )
    return result

class IpTablesController:
    def __init__(self):
        pass

    def execute(self, rule):
        print(f"Executed rule: {rule}")
        return True

    def get_rules(self):
        result = run_shell_cmd("iptables", "-S", "INPUT")
        if result.returncode != 0:
            print(f"Error getting rules {result.stderr}")
        # Here parse rules from cli output
        print(f"stdout rules: {result.stdout.decode('utf-8')}")
        rules = result.stdout.decode('utf-8').split("\n")
        return rules

def main():
    httpKnocker = HttpKnockerPullerClient("http-knocker.altucornet", 8001, "/puller/test")
    ipTables = IpTablesController()

    while True:
        time.sleep(5)
        httpKnocker.push_frw_rules(ipTables.get_rules())
        accepted_rules = []
        rules = httpKnocker.get_last_updates()
        if rules is None or rules == "":
            continue
        print(f"Rules arr size: {len(rules)}")
        for rule in rules:
            # If rule successfully executed then add it's id to accepted list
            if ipTables.execute(rule):
                accepted_rules.append(rule["id"])
        httpKnocker.accept_updates(accepted_rules)



if __name__ == "__main__":
    main()
