from requests import get, post, put, packages
from re import findall
from os.path import exists
import json
import os
packages.urllib3.disable_warnings()


def getcookie(key):
    url = 'https://api.m.jd.com/client.action'
    headers = {
        'cookie': os.environ.get('wsKey'),
        'User-Agent': 'okhttp/3.12.1;jdmall;android;version/10.1.2;build/89743;screen/1440x3007;os/11;network/wifi;',
        'content-type': 'application/x-www-form-urlencoded; charset=UTF-8',
        'charset': 'UTF-8',
        'accept-encoding': 'br,gzip,deflate'
    }
    params = {
        'functionId': 'genToken',
        'client': os.environ.get('client'),
        'clientVersion': os.environ.get('clientVersion'),
        'lang': 'zh_CN',
        'st': os.environ.get('st'),
        'uuid': os.environ.get('uuid'),
        'openudid': os.environ.get('openudid'),
        'sign': os.environ.get('sign'),
        'sv': os.environ.get('sv')
    }
#     print(f"{params}")
#     params = os.environ.get('SIGN')

#     data = 'body=%7B%22to%22%3A%22https%253a%252f%252fplogin.m.jd.com%252fjd-mlogin%252fstatic%252fhtml%252fappjmp_blank.html%22%7D&'
    data = os.environ.get('BODY')
    aa= post(url=url, headers=headers, params=params, data=data, verify=False)
    totokenKey = aa.json()['tokenKey']
    url = 'https://un.m.jd.com/cgi-bin/app/appjmp'



    params = {
        'tokenKey': totokenKey,
        'to': 'https://plogin.m.jd.com/cgi-bin/m/thirdapp_auth_page?token='+totokenKey,
        'client_type': 'android',
        'appid': 879,
        'appup_type': 1,
    }
    res = get(url=url, params=params, verify=False, allow_redirects=False).cookies.get_dict()
    print(f"{res}")
    pt_pin = res['pt_pin']
    cookie = f"pt_key={res['pt_key']};pt_pin={pt_pin};"
    print(f"{res}")
    print(f"{cookie}")
    return pt_pin, cookie

def main():
    cc = os.environ
#     print(f"{os.environ.get('BODY')}")
#     print(f"{os.environ.get('SIGN')}")
#     print(f"111 {os.environ.json()}")
    pin, cookie = getcookie('pin=xxxxx;wskey=xxxxxx;')

if __name__ == '__main__':
    main()
