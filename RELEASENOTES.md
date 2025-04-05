## [v1.1.0]


### Chore


- sort imports to place fyne.io imports above other imports for consistency

- replace 'interface{}' with 'any'

- found more instances of interface{} to replace with any

- remove unused hashBytes func



### Docs


- add inline comments to improve code readability



### Feat


- new countdown timer to next session

- initial support for themes

- initial preferences framework

- add updateTabIfChanged to conditionally refresh UI on hash change

- highlight next race row with primary text color

- add spoiler button to current race to show results



### Fix


- delay setting tray icon for Windows

- load Upcoming Tab on load rather than on click

- correct remnants from old "F1Tray" application

- conditionally hide sprint header on non-sprint events

- remove redundant bulk endpoint load on launch

- headers on schedule table were incorrect

- prevent table cell renderer crash by wrapping canvas text



### Refactor


- export ColoredText widget and constructor




## [v1.0.32]
