import urllib.request
import json

def search_github(query):
    url = f"https://api.github.com/search/repositories?q={urllib.parse.quote(query)}&sort=stars&order=desc"
    req = urllib.request.Request(url, headers={'User-Agent': 'Mozilla/5.0'})
    try:
        with urllib.request.urlopen(req) as response:
            data = json.loads(response.read().decode())
            for item in data.get('items', [])[:5]:
                print(f"- {item['name']} ({item['html_url']}): {item['description']}")
    except Exception as e:
        print(f"Error searching for {query}: {e}")

print("Looking for network traffic generators...")
search_github("network traffic generator")
print("\nLooking for synthetic monitoring...")
search_github("synthetic monitoring")
