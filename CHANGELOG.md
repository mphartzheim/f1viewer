# 🛠️ f1viewer changelog powered by git-cliff

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


### Chore


- clean appimage when building

- still testing cliff changelogs

- still gesting git-cliff out



### Docs


- update changelog for v1.0.32



### Fix


- pass version to git-cliff to enable sectioned changelog

- prevent crash on empty release notes

- correct RELEASENOTES.md extraction using sed instead of awk

- prevent double 'v' in version headings




## [v1.0.31]


### Chore


- treat unreleased changes as next tag with version heading

- fix double v in release notes (again)



### Docs


- update changelog for v1.0.31




## [v1.0.30]


### Chore


- update git-cliff integration

- additional cliff changes

- still testing cliff changes

- still testing cliff

- makefile changes for release

- remove legacy RELEASENOTES.md

- makefile adjustments to check for existing tag



### Docs


- update changelog for v1.0.30

- update changelog for v1.0.30




## [v1.0.27]



## [v1.0.26]



## [v1.0.25]



## [v1.0.24]



## [v1.0.23]



## [v1.0.22]



## [v1.0.21]



## [v1.0.20]



## [v1.0.19]



## [v1.0.18]



## [v1.0.17]



## [v1.0.16]



## [v1.0.15]



## [v1.0.14]



## [v1.0.13]



## [v1.0.12]



## [v1.0.11]



## [v1.0.10]



## [v1.0.9]



## [v1.0.8]



## [v1.0.7]



## [v1.0.6]



## [v1.0.5]



## [v1.0.4]



## [v1.0.3]



## [v1.0.2]



## [v1.0.1]



## [v1.0.0]


