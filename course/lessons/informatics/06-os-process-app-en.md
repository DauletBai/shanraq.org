# Operating system, process, and application

_Lead (summary):_ **We separate a program file, a running process, and OS services to understand launch, permissions, and errors.**

## Where we are on the map

The assistant icon stayed on screen, but its window vanished when you closed it. Did the program itself disappear? Let us break launch into observable steps.

![Operating system, process, and application](/static/course/informatics/map-06-os-process-en.svg)

## A familiar image and the precise model

An icon refers to a program file. At launch, the OS allocates memory, creates a process, and connects it to files and devices. The OS manages processes, memory, files, devices, and access. A driver links a general OS command to hardware. Permission does not prove safety.

## What actually starts

A **program file** is saved instructions on storage; an **application** is a program intended to serve a user's task. A **process** is one running instance with allocated resources; two windows may belong to one process or several. The **operating system (OS)** allocates processor time and memory, opens files, and controls access to devices. A **resource** here is something a process uses: a file, memory, camera, or network. A **permission** is the right to perform a particular action, not evidence that the app is trustworthy.

Act out the assistant's launch with cards labelled `application`, `OS`, `process`, `tasks.json`, and `screen`. The program file rests on storage; the OS checks launch and creates a process; the process asks for `tasks.json`; the OS checks the path and permission; the process displays the result. Remove the `tasks.json` card: the process may start but cannot find data. Remove read permission: the file exists, yet the OS denies access. You have now separated “not found,” “permission denied,” and “program did not start.”

In the project, give each failure three lines: the action that failed, the message observed, and the first safe repair. Giving every permission at once is poor diagnosis: it hides the cause and expands access. The **principle of least privilege** means granting only the right needed for the intended action.

## The limit of the analogy

The familiar image opens the topic but does not replace the mechanism. State where the analogy stops being accurate.

## The lesson's support signal

`program file → OS → process → resource → result`

## Worked example

When opening a file, the OS checks its path and access, chooses an app, creates a process, and passes the file. An error should name the first failed action.

## Predict before observing

Before the experiment, hide the continuation and predict the result in writing. Give the chain of causes from the support signal as well as the answer. Perform the action, record the observation, and compare it with the prediction. A match without an explanation is luck rather than understanding; a mismatch is useful evidence about the link that needs rebuilding.

## Recall without a prompt

Hide the page, rebuild the signal, and explain every transition with a new example.

Explain the topic to someone who has not read the lesson. Do not use a new term until you have unpacked it in plain words. Then restore its precise name and show its place in the system. The listener should be able to predict the next step and give a different example. If they merely repeat a sentence, reduce the explanation to the support signal and rebuild it.

## Find and correct the mistake

“Grant every permission so errors stop.” This hides the cause and increases possible harm. Require least privilege.

## Transfer to a new setting

Apply the model to a device or file absent from the lesson. Separate observation from assumption.

## Project change

List process resources: data file, display, and keyboard. Give each least privilege and a clear denial. Camera, network, and location are unnecessary.

Keep four lines in the decision log: the intended result, the observation, the reason for the chosen action, and the verification. A screenshot can support evidence but cannot replace text and the working file. Do not use personal data: three fictional records provide the same testability and let you show the project safely to another person.

## Exercise

**Required.** Complete the example, correct the error, and make the project change with a reason.

**With your own data.** Test the decision on three fictional records.

**Optional.** Ask another person to rebuild the explanation from the map.

The readiness criterion is concrete: reproduce the signal without the page, correct the proposed error, transfer the model to an unfamiliar example, and show the project change. If one part fails, return to that link instead of rereading the whole lesson. After repairing it, repeat only an equivalent task with different data.

## Retrieval after 1, 7, and 30 days

Rebuild the signal tomorrow; solve an equivalent task after seven days; after thirty days, verify the decision in the project.

[Next lesson: files, folders, and exact paths](/read/informatics-07-files-folders-paths?lang=en)
