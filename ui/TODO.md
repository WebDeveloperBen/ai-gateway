- [x] add teams page
- [x] add roles page

# Global

- [ ] create DRY reusable composable or something to handle permission checks across the application at a route level and then a component level

# Backend Api

- [ ] add the default roles into the database and properly define their permissions
- [ ] add published status to policies

# Applications page

- [ ] improve the ApiKeysReveal.vue component, warn and ask for confirmation from the user before regnerating the api key. it's a destructive action

# Environment page

- [ ] add action to team member add member card
- [ ] add action to applications deploy app card
- [ ] add action when hovering and clicking api keys
- [ ] add refresh action to the recent activity card icon button
- [ ] add modal to edit the environment

# Analytics page

- [ ] Add filtering of charts to date ranges
- [ ] allow filtering of charts by applications and environments
- [ ] allow filtering using url params and deep linking
- [ ] add log driven view for raw span trace logging outputting (pagination is a must)

# Governance page

- [ ] setup alerts page
- [ ] setup policies page
- [ ] setup monaco editor
- [ ] add CEL support to monaco editor for syntax highlighting and linting

# Prompts Page

- [ ] save prompts and retrieve by id
- [ ] move save button to end of row
- [ ] move ready and its icon to somewhere else or even just remove it
- [ ] set prompts per environment
- [ ] set prompts per app
- [ ] versioned prompts
- [ ] promote to live prompt in use by apps
- [ ] fix the load prompts button so it actually works
- [ ] break up the massive page into DRY components
- [ ] make the model paramaters actually do something
- [ ] create backend endpoint to call prompt playground model (should be the account holders model not ours)
- [ ] create ability to save the settings of the paramaters to the prompt
- [ ] create logic to fetch latest version prompt from the backend
- [ ] create ability to let users define which versioning schema they want to use (Nice to have only)
- [ ] add prompt templates to database
- [ ]

# Prompt tags

- [ ] Let users tag their prompts to enhance filtering, organisation, and improve user interface styling, user defined through strings, no impact in the application just for user experience purpose and should feed into filtering

# Budgets

- [ ] How to maintain a cost sheet with live data ? sync daily? from each provider that we support ? -> required for budget policies

# Policy page

- [ ] make the validate button do something
- [ ] fix the hardcoded validating status
- [ ] make the ai draft do something
- [ ] add keyboard short cuts to pop open a ai draft box to get ai help
- [ ] handle the idea of published and non published draft policies
- [ ] handle validate and deploy functionality
- [ ] expand the monaco editor so it works with treesitter
- [ ] expand the monaco editor so it works with linting
- [ ] expand the monaco editor so it doesn't allow duplicate expressions
- [ ] create all the default expressions for the templates
- [ ] load the policy templates via the api
- [ ] Move the policy page into the governance section
- [ ] create policy / id page or modal
- [ ] Move the policy creation page to an internal page

# Providers

- [ ] add the ability to add custom models
- [ ] add ability to add fine tuned models

# Research

When it comes to do sass use this:

- [ ] <https://ui-thing.behonbaker.com/components/dialog#change-plan>
- [ ] add <https://ui-thing.behonbaker.com/components/dialog#onboarding>
- [ ] add <https://ui-thing.behonbaker.com/components/dialog#edit-profile>
- [ ] add <https://ui-thing.behonbaker.com/goodies/terminal#usage>
- [ ] add <https://ui-thing.behonbaker.com/goodies/tip-tap> (or reuse monaco editor)
- [ ] explore replacing native select with this: <https://ui-thing.behonbaker.com/components/listbox#objects>
