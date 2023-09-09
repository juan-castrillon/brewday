# BrewDay


![logo](https://github.com/juan-castrillon/brewday/assets/64461123/5c7ad1bd-34d9-4b3f-8621-0c97e55d77a5)

BrewDay is a self-contained web application aimed at helping homebrewers with their brewing process. 

It is intended to be used the day of the brew, and it is designed to be used on all devices, from desktop to mobile.

The app is intended to be self-hosted and does not have multiple users in mind. It is designed to be used by a single user at a time.

The app helps the user with the following tasks:

- **Follow the recipe**. The user can import a recipe from any of the supported formats (see below), and the app will guide the user through the brewing process, step by step. 
- **Note taking**. The user can take notes during the brew, and the app will save them for future reference. Each step in the process gives the opportunity to input real data (to compare with the recipe) and notes (to keep track of the brew).
- **Timers**. The app will set timers for each step in the process, and will notify the user when the time is up. - **Statistics**. The app will calculate the efficiency of the brew, evaporation rate, and other useful statistics.
- **Timeline and summary**. The app will ley the users download a timeline of the brew, and a summary of the brew day, with all the relevant data. Supported summary formats are listed below.

## Supported recipe formats

The app supports the following recipe formats:
- [Maische Malz und Mehr](https://www.maischemalzundmehr.de/index.php?inhaltmitte=lr) ([JSON](https://www.maischemalzundmehr.de/rezept.json.txt))
- [Braureka](https://braureka.de/) (JSON) (This is supposed to be MMUM, but it differs in implementation of some fields that are parsed as strings instead of numbers)


## Supported summary formats

The app supports the following summary formats:
- [Markdown](https://www.markdownguide.org/basic-syntax/): Markdown summary will create a summary of the brew day in Markdown format. This is useful to copy and paste the summary in a blog post, or to share it with other people. The timeline is just a list of timestamps. 
