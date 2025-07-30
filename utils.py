from dotenv import load_dotenv
from os import getenv
from google.oauth2.service_account import Credentials
from googleapiclient.discovery import build
import json, asyncio

load_dotenv()

SPREADSHEET_ID = getenv("SPREADSHEET_ID")
SCOPES = json.loads(getenv("SCOPES"))

credentials = Credentials.from_service_account_file("key.json", scopes=SCOPES)
service = build("sheets", "v4", credentials=credentials)

names_hours_list = []

def url_to_id(url):
    try:
        return url.split("d/")[1].split("/edit")[0]
    except IndexError:
        try:
            return url.split("document_id=")[1]
        except IndexError:
            return url

async def update_hours_list(names_col_arg, nicknames_col_arg, year_col_arg, term_hour_col_arg, all_hours_call_arg):
    global names_hours_list
    names_hours_list.clear()

    names_hours_data_request = await asyncio.to_thread(
        service.spreadsheets().values().batchGet,
        spreadsheetId=SPREADSHEET_ID,
        ranges=[names_col_arg, nicknames_col_arg, year_col_arg, term_hour_col_arg, all_hours_call_arg]
    )
    names_hours_data = names_hours_data_request.execute()

    names_len = len(names_hours_data["valueRanges"][0]["values"])
    nicknames_len = len(names_hours_data["valueRanges"][1]["values"])

    for i in range(names_len):
        last, first = names_hours_data["valueRanges"][0]["values"][i][0].split(", ")

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

def get_hours(name):
    global names_hours_list

    if len(names_hours_list) == 0:
        return None

    name = name.lower()

    for value in names_hours_list:
        if name in value["name"] or name in value["nickname"]:
            return value
    return None

def get_year_ranking(year):
    global names_hours_list

    if len(names_hours_list) == 0:
        return None

    year = year.lower()
    year_ranking = [each for each in names_hours_list if each["year"] == year]

    for i in range(len(year_ranking)):
        for j in range(i + 1, len(year_ranking)):
            if year_ranking[i]["all_hours"] < year_ranking[j]["all_hours"]:
                year_ranking[i], year_ranking[j] = year_ranking[j], year_ranking[i]

    return year_ranking[0:5]

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