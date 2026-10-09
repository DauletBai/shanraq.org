"""Command-line Digital Assistant 1.0. Run: python3 assistant.py --help."""

import argparse
from datetime import date
from pathlib import Path

from assistant_core import (add_task, complete_task, count_done, load_document,
                            parse_date, reminder_status, save_document)


def main(argv=None):
    parser = argparse.ArgumentParser(description="My Digital Assistant 1.0")
    parser.add_argument("--file", type=Path, default=Path(__file__).with_name("tasks.json"),
                        help="UTF-8 JSON file with fictional tasks")
    commands = parser.add_subparsers(dest="command", required=True)
    commands.add_parser("list", help="show all tasks and a completion count")
    reminders = commands.add_parser("reminders", help="show each task's reminder status")
    reminders.add_argument("--today", help="YYYY-MM-DD; defaults to local calendar date")
    add = commands.add_parser("add", help="add a fictional task")
    add.add_argument("id")
    add.add_argument("title")
    add.add_argument("--due", help="YYYY-MM-DD")
    done = commands.add_parser("done", help="mark one task complete")
    done.add_argument("id")
    args = parser.parse_args(argv)

    try:
        document = load_document(args.file)
        if args.command == "list":
            for task in document["tasks"]:
                mark = "x" if task["done"] else " "
                print(f"[{mark}] {task['id']} {task['title']}")
            print(f"Completed: {count_done(document['tasks'])}/{len(document['tasks'])}")
        elif args.command == "reminders":
            today = parse_date(args.today) if args.today else date.today()
            print(f"Today: {today.isoformat()} (local calendar date)")
            for task in document["tasks"]:
                print(f"{task['id']}: {reminder_status(task, today)}")
        elif args.command == "add":
            updated = add_task(document, args.id, args.title, args.due)
            save_document(args.file, updated)
            print(f"Added: {args.id}")
        elif args.command == "done":
            updated = complete_task(document, args.id)
            save_document(args.file, updated)
            print(f"Completed: {args.id}")
    except ValueError as exc:
        parser.exit(2, f"Error: {exc}\n")


if __name__ == "__main__":
    main()
