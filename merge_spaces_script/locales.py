
import json
import csv
import sys
import json

# Downloads Contenful Locales json files per space provided
#python -c 'import json; fp = open("migrations/locales_test.json", "r"); obj = json.load(fp); fp.close(); for restaurant in data["items"]: print restaurant["internal_code"]'

spaces_arg = sys.argv[1]

spaces = spaces_arg.split(",")
langs = []
langs_json = []

for space in spaces: 
    with open('migrations/locales_'+space+'.json') as f:
        data = json.load(f)
        for item in data["items"]: 
            lang_code = item["internal_code"]
            lang_name = item["name"]
            if lang_code not in langs:
                x = {
                  "name": lang_name,
                  "code": lang_code
                }
                # convert into JSON:
                y = json.dumps(x)
                langs_json.append(y)
                langs.append(lang_code)

rows_file = open('migrations/all_spaces_locales.txt', 'w') #write to file
rows_file.writelines("%s\n" % i for i in langs_json)
rows_file.close()