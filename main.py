from discord import Intents, Activity, ActivityType, Embed, Color, Interaction, app_commands
from discord.ext import commands
from google.oauth2.service_account import Credentials
from googleapiclient.discovery import build

from dotenv import load_dotenv
from os import getenv
from time import time
from typing import Optional
import json

from utils import update_hours_list, get_hours, find_default_name, write_default_name


load_dotenv()

token = getenv("DISCORD_TOKEN")
scopes = json.loads(getenv("SCOPES"))

update_delay = int(getenv("update_delay"))

names_hours_list = []
credentials = Credentials.from_service_account_file("key.json", scopes=scopes)
sheets_service = build("sheets", "v4", credentials=credentials)

intents = Intents.default()
intents.message_content = True
client = commands.Bot(command_prefix="/", intents=intents)
last_update = 0

@client.event
async def on_ready():
    global names_hours_list, sheets_service
    print(f"Logged in as {client.user}")
    await client.tree.sync()
    await client.change_presence(activity=Activity(type=ActivityType.playing, name="Watching your messages"))
    await update_hours_list(names_hours_list, sheets_service)

@client.tree.command(name="hours", description="Returns the hours of a given name or default name")
@app_commands.describe(name="Returns the hours of a given name or default name")
async def hours(interaction: Interaction, name: Optional[str] = None):
    global names_hours_list
    # checks for default name, asks to set if not found
    if not name:
        if name_found := find_default_name(interaction.user.id):
            name = name_found
        else:
            # embed = Embed(
            #     title="Default name not set",
            #     description="You don't have a default name set, set it using `/setname <name>`",
            #     color=Color.pink()
            # )
            # await interaction.response.send_message(embed=embed, ephemeral=True)
            await interaction.response.send_message("https://tenor.com/view/but-none-were-there-spongebob-hawaii-part-ii-gif-1736518424783391175", ephemeral=False)
            return

    # retrieves and sends hours info, send error if not found
    if hours_info := get_hours(names_hours_list, name):
        embed = Embed(
            title=f"{hours_info['name']}'s Hours",
            description=f"Term Hours: {hours_info['term_hours']}\nAll Hours: {hours_info['all_hours']}",
            color=Color.gold()
        )
        await interaction.response.send_message(embed=embed, ephemeral=False)
    else:
        # embed = Embed(
        #     title="Name not found",
        #     description="The name you provided was not found in the spreadsheet",
        #     color=Color.pink()
        # )
        # await interaction.response.send_message(embed=embed, ephemeral=True)
        await interaction.response.send_message("https://tenor.com/view/but-none-were-there-spongebob-hawaii-part-ii-gif-1736518424783391175", ephemeral=False)
        return

    # if name was given but it's not in the default names, ask to set it
    if not find_default_name(interaction.user.id):
        embed = Embed(
            title="Default name not set",
            description="You don't have a default name, set it to make this command easier",
            color=Color.pink()
        )
        await interaction.response.send_message(embed=embed, ephemeral=True)

@client.tree.command(name="setname", description="Sets a default name for the /hours command")
@app_commands.describe(name="Sets the default name for the /hours command")
async def setname(interaction: Interaction, name: str):
    write_default_name(interaction.user.id, name)
    await interaction.response.send_message(f"Default name set to {name}", ephemeral=True)

@client.tree.command(name="updatehours", description="Updates the /hours command")
async def updatehours(interaction: Interaction):
    global last_update

    if time() - last_update <= update_delay:
        await interaction.response.send_message(
            f"Update request is too soon, you must wait {int(update_delay - (time() - last_update))} more seconds to update the bot",
            ephemeral=True)
        return

    last_update = time()
    await update_hours_list(names_hours_list, sheets_service)
    await interaction.response.send_message(
        f"Hours have been updated, you may update them again in 5 minutes",
        ephemeral=False)

@client.tree.command(name="help", description="Shows the help message")
async def help(interaction: Interaction):
    embed = Embed(
        title="Help",
        description="""
Here are this bots commands:
- `/help`: Shows this message
- `/hours`: Shows your hours if you have a set name
- `/hours {name}`: Shows the hours of the person with the name specified
- `/setname {name}`: Sets the default name for you with the name specified
- `/updatehours`: Updates the hours for the bot, 5 min cooldown
        """,
        color=Color.gold()
    )
    await interaction.response.send_message(embed=embed, ephemeral=True)

client.run(token)
