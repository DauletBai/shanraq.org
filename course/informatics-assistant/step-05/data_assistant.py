"""Create and inspect the fictional SQLite checkpoint; no external packages."""

import argparse
from pathlib import Path

from data_store import create_database, minutes_report


ROOT = Path(__file__).resolve().parent


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("command", choices=("init", "report"))
    parser.add_argument("--db", type=Path, default=ROOT / "assistant.db")
    args = parser.parse_args()
    if args.command == "init":
        create_database(args.db, ROOT / "tasks.json", ROOT / "study_sessions.csv")
        print("SQLite checkpoint created: 3 tasks, 4 fictional study sessions")
    else:
        for task_id, minutes in minutes_report(args.db):
            print(task_id, minutes)


if __name__ == "__main__":
    main()
