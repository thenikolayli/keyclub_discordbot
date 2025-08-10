# collection of utility functions used to automatically log events and meetings and check volunteer hours

from dotenv import load_dotenv
from os import getenv
import asyncio, json

load_dotenv()

spreadsheet_id = getenv("SPREADSHEET_ID")
names_col = getenv("NAMES_COL")
nicknames_col = getenv("NICKNAMES_COL")
year_col = getenv("YEAR_COL")
term_hours_col = getenv("TERM_HOURS_COL")
all_hours_col = getenv("ALL_HOURS_COL")
spreadsheet_ranges = [names_col, nicknames_col]


# ------------------CHECK HOURS API STUFF------------------
# this is the same as in the website api, it's just independent

# updates the hours list by fetching hours from the spreadsheet
async def update_hours_list(names_hours_list, service):
    global names_col, nicknames_col, year_col, term_hours_col, all_hours_col, spreadsheet_id

    names_hours_list.clear()

    names_hours_data_request = await asyncio.to_thread(
        service.spreadsheets().values().batchGet,
        spreadsheetId=spreadsheet_id,
        ranges=[names_col, nicknames_col, year_col, term_hours_col, all_hours_col]
    )
    names_hours_data = names_hours_data_request.execute()
    nicknames_len = len(names_hours_data["valueRanges"][1]["values"])
    loop_range = len(names_hours_data["valueRanges"][2]["values"])

    for i in range(loop_range):
        last, first = names_hours_data["valueRanges"][0]["values"][i][0].split(", ")
        # print(first, last, i)

        full_name = f"{first.lower()} {last.lower()}"
        if i >= nicknames_len or names_hours_data["valueRanges"][1]["values"][i] == []:
            nickname = ""
        else:
            nickname = names_hours_data["valueRanges"][1]["values"][i][0].lower()
        year = names_hours_data["valueRanges"][2]["values"][i][0].lower()
        term_hours = float(names_hours_data["valueRanges"][3]["values"][i][0])
        all_hours = float(names_hours_data["valueRanges"][4]["values"][i][0])

        names_hours_list.append({
            "name": full_name,
            "nickname": nickname,
            "year": year,
            "term_hours": term_hours,
            "all_hours": all_hours
        })

# gets the hours for a person based on their name
def get_hours(names_hours_list, name):
    if len(names_hours_list) == 0:
        return None

    name = name.lower()

    for value in names_hours_list:
        if name in value["name"] or name in value["nickname"]:
            return value
    return None


# ------------------DEFAULT NAMES API STUFF------------------

def find_default_name(user_id):
    with open("default_names.json", "r") as file:
        default_names = json.load(file)
    print(default_names.get(str(user_id)))
    return default_names.get(str(user_id), None)

def write_default_name(user_id, name):
    with open("default_names.json", "r") as file:
        default_names = json.load(file)

    default_names[str(user_id)] = name

    with open("default_names.json", "w") as file:
        json.dump(default_names, file)