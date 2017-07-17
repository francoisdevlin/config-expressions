The enum matcher is a valuable way to use a union type.  However, you must explicitly include every match you want to include.  Sometimes it is more useful to match a more general pattern, such as a regular expression.  This example shows a regex matcher at work.  The regex is specified in the expression `/app\\w/$app_name.db.user`. Notice the following

* The regex is specified with `/` delimters on each end.  
* The exact value of the regex is being captured in the variable `app_name`
* The regex is automatically anchored.
* The regex has high precedence than a wildcard
