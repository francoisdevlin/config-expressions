This example shows to to use an enum, with variable substitution.  You'll notice the following:

* The key `app1,app2,app3$app_name.db.user` specifies an enumeration for 3 apps
* The exact value of the enumeration is being captured in the variable `app_name`
* The enum has high precedence than a wildcard
