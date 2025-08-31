# Plan

## Tentative structure

### Go CLI to upsert notices

#### Input
1. File names 
2. Config options
    1. company_name (eg: Acme Inc)
    2. notice_format (eg: "Copyright (C) $company_name $year")
    3. auto_fix (boolean flag)

#### Output
1. Raise error code to fail pre-commit if any files fail check

### Shell script

This shell script will hook into the pre-commit event in git and
call our CLI with the file names that were part of changed files


## Milestones 

1. Create a GO CLI
    1. Takes in filenames 
    2. Check copyright notices, raise error code if need
    3. print filenames without copyright notices

2. Feature to add copyright notices
    1. It should only add notices wherever missing

3. Feature to update copyright notices
    1. Add boolean flag to auto_fix
    2. Should print number of files updated

3. Explore integration with [pre-commit framework](https://pre-commit.com/#golang)

5. Feature for custom template for copyright notices 
