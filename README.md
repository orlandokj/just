# Just

Just is the CLI application to run projects that you don't want to work on. 
You just need to run it, and thats where the name came from.


## Motivation

At my job I had several projects that I need to run to test some features as a whole.
To run that amount of projects on dev mode is really expensive to my poor laptop.
So I created this project in order to standardize the way I run these projects.


## Type of projects
- Static file server, you need to define name, port, and the directory that 
will be served
- Server, you need to define the command to run the application.

## Steps
There are (for now) two main steps to run these kind of projects:
1. Build: If you want Just to build the project for you, you need to define how
to build it, defining a command to build it. Keep in mind that you are building 
to run on your local machine, you probrably need to pass somekind of profile to 
your build to use your developer settings.
2. Run: This is the default action and the main reason for this tool, indicate 
how to run your project.

TODO More detailed about configuration
