from requests import get, post, put, packages
from re import findall
from os.path import exists
import json
import os
packages.urllib3.disable_warnings()


def getcookie():
    res = get(url="https://api.jds.codes/gentoken", verify=False, allow_redirects=False).json()["data"]
    url = 'https://api.m.jd.com/client.action'
    headers = {
        'cookie': os.environ.get('wsKey'),
        'User-Agent': 'okhttp/3.12.1;jdmall;android;version/10.1.2;build/89743;screen/1440x3007;os/11;network/wifi;',
        'content-type': 'application/x-www-form-urlencoded; charset=UTF-8',
        'charset': 'UTF-8',
        'accept-encoding': 'br,gzip,deflate'
    }
    url = 'https://api.m.jd.com/client.action?functionId=genToken&'+res['sign']
    aa = post(url=url, headers=headers, verify=False)
    url = 'https://un.m.jd.com/cgi-bin/app/appjmp'
    params = {
        'tokenKey': aa.json()['tokenKey'],
        'to': 'https://plogin.m.jd.com/cgi-bin/m/thirdapp_auth_page?token='+aa.json()['tokenKey'],
        'client_type': 'android',
        'appid': 879,
        'appup_type': 1,
    }
    res = get(url=url, params=params, verify=False, allow_redirects=False).cookies.get_dict()
    cookie = f"pt_key={res['pt_key']};pt_pin={res['pt_pin']};"
    print(f"{cookie}")
    return params

def main():
      getcookie()

if __name__ == '__main__':
    main()
