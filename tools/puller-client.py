import re
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
            data={"accepted_commands": json.dumps(accepted_rule_ids)},
        )
        print(f"accepted response: {r.text}")

    def push_frw_rules(self, rules_set):
        print(f"rule set count: {len(rules_set)}")
        print(f"rule set: {rules_set}")
        rules_body_dict = []
        for rule in rules_set:
            rule.debug()
            rules_body_dict.append(rule.get())
        r = requests.post(
            self.__base + "/pushRulesSet", data={"rules": json.dumps(rules_body_dict)}
        )
        print(f"rule set reponse: {r.text}")


def run_shell_cmd(*args):
    result = subprocess.run(args, capture_output=True)
    return result


class IpTablesRule:
    def __init__(self):
        self.__body = {}
        self.__keys = [
            "id",
            "action",
            "chain",
            "disabled",
            "protocol",
            "src-address",
            "dst-port",
            "comment",
            "place-before",
        ]
        self.__re = {}
        self.__re["chain"] = re.compile("-A\s([^\s]+)\s")
        self.__re["action"] = re.compile("\s-j\s([^\s]+)")
        self.__re["protocol"] = re.compile("\s-p\s([A-Za-z]+)")
        self.__re["src-address"] = re.compile("\s-s\s([^\s]+)")
        self.__re["dst-port"] = re.compile("\s--dport\s([^\s]+)")
        self.__re["comment"] = re.compile(
            "-m\s+comment\s+--comment\s+(\"[^\"]*\"|'[^']*'|[^'\"\s]+)"
        )

    def setId(self, id):
        self.__body["id"] = id

    def isEmpty(self):
        return len(self.__body) == 0

    def get(self):
        return self.__body

    def debug(self):
        print(f"iptables rule dbg: {self.__body}")

    def __extract_regex_by_key(self, key, input):
        if key in self.__re:
            match = self.__re[key].search(input)
            # print(f"match for {key} => {match}")
            if match != None and len(match.groups()) != 0:
                self.__body[key] = match.groups()[0]

    def from_string(self, line):
        for key in self.__keys:
            self.__extract_regex_by_key(key, line)

    def from_dict(self, dict):
        for key in self.__keys:
            self.__body[key] = dict.get(key)

    def get_shell_args(self):
        args = []
        args.append(f"{self.__body['chain'].upper()}")
        # 0 in place-before means that drop rule not found, at least for now...
        if (
            self.__body["place-before"] is not None
            and self.__body["place-before"] != "0"
        ):
            args.append(f"{self.__body['place-before']}")
        args.append("-p")
        args.append(f"{self.__body['protocol']}")
        args.append("-s")
        args.append(f"{self.__body['src-address']}")
        args.append("--dport")
        args.append(f"{self.__body['dst-port']}")
        args.append("-j")
        args.append(f"{self.__body['action'].upper()}")
        args.append("-m")
        args.append("comment")
        args.append("--comment")
        args.append(f"{self.__body['comment']}")
        return args


class IpTablesController:
    def __init__(self):
        pass

    def execute(self, command):
        print(f"Executed command: {command}")
        if command["command"]["type"] == "remove":
            return self.execute_remove(command)
        elif command["command"]["type"] == "add":
            return self.execute_insert(command)
        return False

    def execute_remove(self, command):
        # iptables --delete INPUT 3
        # +1 in rule index because numeration starts from 1
        # BUT
        # -P INPUT ACCEPT -- this rule can be treated as 0 index
        args = [
            "iptables",
            "--delete",
            "INPUT",
            f'{command["command"]["id"]}',
        ]  # -I - insert
        result = run_shell_cmd(*args)
        if result.returncode != 0:
            print(f"Error exeuting REMOVE rule {result.stderr}")
            return False
        return True

    def execute_insert(self, command):
        rule = IpTablesRule()
        rule.from_dict(command["command"]["rule"])
        args = ["iptables", "-I"]  # -I - insert
        args.extend(rule.get_shell_args())
        result = run_shell_cmd(*args)
        if result.returncode != 0:
            print(f"Error executing INSERT rule {result.stderr}")
            return False
        return True

    def get_rules(self):
        result = run_shell_cmd("iptables", "-S", "INPUT")
        if result.returncode != 0:
            print(f"Error getting rules {result.stderr}")
        # Here parse rules from cli output
        print(f"stdout rules: {result.stdout.decode('utf-8')}")
        rule_lines = result.stdout.decode("utf-8").split("\n")
        rules = []
        index = 0
        for line in rule_lines:
            if line == "":
                continue
            print(f"making rule from line: {line}")
            rule = IpTablesRule()
            rule.from_string(line)
            # if rule.isEmpty():
            #     continue
            rule.setId(index)
            rules.append(rule)
            index += 1
        return rules


def main():
    httpKnocker = HttpKnockerPullerClient("http-knocker.skyline", 8001, "/puller/test")
    ipTables = IpTablesController()

    while True:
        time.sleep(5)
        httpKnocker.push_frw_rules(ipTables.get_rules())
        accepted_rules = []
        commands = httpKnocker.get_last_updates()
        if commands is None or commands == "":
            continue
        print(f"Commands arr size: {len(commands)}")
        for command in commands:
            # If rule successfully executed then add it's id to accepted list
            if ipTables.execute(command):
                accepted_rules.append(command["id"])
        httpKnocker.accept_updates(accepted_rules)


def test_parsing():
    line = "-A INPUT -s 127.0.0.1/32 -p tcp -m tcp --dport 3333 -m comment --comment httpKnocker-basicfirewall-1680912119-dc7fe68f27ff5f2a9fd71b9f06f6e3dd43ea919a -j ACCEPT"
    rule = IpTablesRule()
    rule.from_string(line)
    rule.debug()


if __name__ == "__main__":
    main()
    # test_parsing()


"""

-P INPUT ACCEPT
-A INPUT -s 127.0.0.1/32 -p tcp -m tcp --dport 3333 -m comment --comment httpKnocker-basicfirewall-1680912119-dc7fe68f27ff5f2a9fd71b9f06f6e3dd43ea919a -j ACCEPT
-A INPUT -s 0.0.0.0/32 -p tcp -m tcp --dport 2222 -m comment --comment http-knocker-drop-all-rule -j DROP

"""
