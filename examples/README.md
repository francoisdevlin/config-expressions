The goal of this document is to show both the _why_ as much as the _how_ to use configuration expressions.  It starts by walking you through an example of a static, vanilla JSON file, and adds more an more features to add more an more power.  At the end is a reference section on each expression

# Learning by Example
Let's start with a very basic example, using vanilla JSON
## 001 vanilla json
Vanilla JSON is a configuration format we are all used to using.  There is a lot of pre-existing configuration written as vanilla JSON, so reading it directly made sense as a starting point for the tool.  It is where you will start when trying to abstract out expressions from pre-existing files


    {
    	"development": {
    		"app1.db.user" : "app1_user",
    		"app1.db.password" : "secret_1",
    		"app1.db.url" : "jdbc:h2:/sample/path",
    
    
    		"app2.db.user" : "app2_user",
    		"app2.db.password" : "secret_2",
    		"app2.db.url" : "jdbc:h2:/sample/path",
    
    
    		"app3.db.user" : "app3_user",
    		"app3.db.password" : "secret_3",
    		"app3.db.url" : "jdbc:h2:/sample/path"
    	}
    }

This will produce the following output

    $ go/confexpr-go lookup development.app1.db.user
    'app1_user'
     
    $ go/confexpr-go lookup development.app2.db.user
    'app2_user'
     
## 002 basic wildcard
Let's us a wildcard pattern to exact out some of the common functionality.  You'll notice the app1.db.url and app3.db.url are all the same.  We can match all of them with a widcard pattern, `*.db.url`.  This will let us remove a couple of lines from the config


    {
    	"development": {
    		"app1.db.user" : "app1_user",
    		"app1.db.password" : "secret_1",
    
    		"app2.db.user" : "app2_user",
    		"app2.db.password" : "secret_2",
    		"app2.db.url" : "jdbc:h2:/other/path",
    
    		"app3.db.user" : "app3_user",
    		"app3.db.password" : "secret_3",
    
    		"*.db.url" : "jdbc:h2:/sample/path"
    	}
    }

This will produce the following output

    $ go/confexpr-go lookup development.app1.db.user
    'app1_user'
     
    $ go/confexpr-go lookup development.app2.db.user
    'app2_user'
     
    $ go/confexpr-go lookup development.app1.db.url
    'jdbc:h2:/sample/path'
     
You can see that `app2.db.url` matches the more specific input

    $ go/confexpr-go lookup development.app2.db.url
    'jdbc:h2:/other/path'
     
We can also get a value for entries that aren't specified explicitly, such as `app9.db.url`

    $ go/confexpr-go lookup development.app9.db.url
    'jdbc:h2:/sample/path'
     
## 003 basic substitution
This example shows to to use basic variable substitution with a wildcard.  It matches keys that are static, such as app1.db, as well as keys that are dynamic, such as app9.


    {
    	"development": {
    		"app1.db.password" : "secret_1",
    		"app2.db.password" : "secret_2",
    		"app3.db.password" : "secret_3",
    
    		"*$app_name.db.user" : "${app_name}_user",
    		"*.db.url" : "jdbc:h2:/sample/path"
    	}
    }

This will produce the following output

Notice that the matched part of the path is substituted into the returned value

    $ go/confexpr-go lookup development.app1.db.user
    'app1_user'
     
The url does not have any substitution, so no changes are made

    $ go/confexpr-go lookup development.app1.db.url
    'jdbc:h2:/sample/path'
     
## 004 basic enum
This example shows to to use an enum, with variable substitution.  You'll notice the following:

* The key `app1,app2,app3$app_name.db.user` specifies an enumeration for 3 apps
* The exact value of the enumeration is being captured in the variable `app_name`
* The enum has high precedence than a wildcard


    {
    	"development": {
    		"app1.db.password" : "secret_1",
    		"app2.db.password" : "secret_2",
    		"app3.db.password" : "secret_3",
    
    		"app1,app2,app3$app_name.db.user" : "enum_${app_name}_user",
    		"*$app_name.db.user" : "${app_name}_user",
    		"*.db.url" : "jdbc:h2:/sample/path"
    	}
    }

This will produce the following output

    $ go/confexpr-go lookup development.app1.db.user
    'enum_app1_user'
     
    $ go/confexpr-go lookup development.app2.db.user
    'enum_app2_user'
     
Notice that the wildcard pattern is matched after the enum is exhausted

    $ go/confexpr-go lookup development.app9.db.user
    'app9_user'
     
## 005 basic regex
The enum matcher is a valuable way to use a union type.  However, you must explicitly include every match you want to include.  Sometimes it is more useful to match a more general pattern, such as a regular expression.  This example shows a regex matcher at work.  The regex is specified in the expression `/app\\w/$app_name.db.user`. Notice the following

* The regex is specified with `/` delimters on each end.  
* The exact value of the regex is being captured in the variable `app_name`
* The regex is automatically anchored.
* The regex has high precedence than a wildcard


    {
    	"development": {
    		"app1.db.password" : "secret_1",
    		"app2.db.password" : "secret_2",
    		"app3.db.password" : "secret_3",
    
    		"/app\\w/$app_name.db.user" : "regex_${app_name}_user",
    		"*$app_name.db.user" : "${app_name}_user",
    		"*.db.url" : "jdbc:h2:/sample/path"
    	}
    }

This will produce the following output

    $ go/confexpr-go lookup development.app1.db.user
    'regex_app1_user'
     
    $ go/confexpr-go lookup development.app2.db.user
    'regex_app2_user'
     
Notice the built in anchors to the regex, so this value falls through to the wildcard pattern

    $ go/confexpr-go lookup development.app10.db.user
    'app10_user'
     
## 006 complete precedence
This example shows the complete precedence of expressions of the same locality.  The order is

1. Direct Match
2. Enum Match
3. Regex Match
4. Wildcard Match
5. Deep Wildcard Match



    {
    	"development": {
    		"app1.db.user" : "user",
    		"app1,app2$app_name.db.user" : "enum_${app_name}_user",
    		"/app\\w/$app_name.db.user" : "regex_${app_name}_user",
    		"*$app_name.db.user" : "wildcard_${app_name}_user",
    		"**$app_name.user" : "deep_wildcard_${app_name}_user"
    	}
    }

This will produce the following output

    $ go/confexpr-go lookup development.app1.db.user
    'user'
     
    $ go/confexpr-go lookup development.app2.db.user
    'enum_app2_user'
     
    $ go/confexpr-go lookup development.app3.db.user
    'regex_app3_user'
     
    $ go/confexpr-go lookup development.app10.db.user
    'wildcard_app10_user'
     
Notice that the entire matched path is substituted

    $ go/confexpr-go lookup development.other_app.server.user
    'deep_wildcard_other_app.server_user'
     
## 007 locality vs specificity
Locality is established by creating a new dictionary object

This example shows local instructions winning oevr global ones.  Observe development.app1.db.url.  Even thought the global expression is more specific, the specification of a local wildcard expression overrides the global one.  This is because the most local expression wins a conflict. 

You can see that more specific expressions still win if they have the same locality, in `development.app2.db.url`


    {
    	"development.app1.db.url" : "jdbc:h2:/not/used",
    	"development.app1.service.url" : "jdbc:h2:/used",
    	"development.app2.db.url" : "jdbc:h2:/also/not/used",
    
    	"development": {
    		"app2.db.url" : "jdbc:h2:/other/path",
    		"*.db.url" : "jdbc:h2:/sample/path"
    	}
    }

This will produce the following output

    $ go/confexpr-go lookup development.app1.db.url
    'jdbc:h2:/sample/path'
     
    $ go/confexpr-go lookup development.app2.db.url
    'jdbc:h2:/other/path'
     
Noting matches locally inside the `development` dictionary, so the engine falls back on the global value

    $ go/confexpr-go lookup development.app1.service.url
    'jdbc:h2:/used'
     
## 008 explain
Sometimes it can be difficult to reason about how each configuation expression is being evaluated.  Fortunately, there is the `explain` command to assit with that.


    {
    	"development.app1.db.url" : "jdbc:h2:/not/used",
    	"development.app1.service.url" : "jdbc:h2:/used",
    	"development.app2.db.url" : "jdbc:h2:/also/not/used",
    
    	"development": {
    		"app2.db.url" : "jdbc:h2:/other/path",
    		"*.db.url" : "jdbc:h2:/sample/path"
    	}
    }

The explain command can be used to gain insight into why the tool returned a certain value for an expression

Here you can see how the engine determines which expression to use for `development.app1.db.url`

    $ go/confexpr-go explain development.app1.db.url
    The following rules were evaluated in this order, the first hit is returned
    Rule     1: MISS   : app2.db.url, ignoring children
    Rule     2: HIT    : *.db.url VALUE: 'jdbc:h2:/sample/path'
    Rule     3: HIT    : development.app1.db.url VALUE: 'jdbc:h2:/not/used'
    Rule     4: MISS   : development.app1.service.url, ignoring children
    Rule     5: MISS   : development.app2.db.url, ignoring children
And here you can see there are no hits locally, and the global value is the highest priority hit

    $ go/confexpr-go explain development.app1.service.url jdbc:h2:/used
    The following rules were evaluated in this order, the first hit is returned
    Rule     1: MISS   : app2.db.url, ignoring children
    Rule     2: MISS   : *.db.url, ignoring children
    Rule     3: MISS   : development.app1.db.url, ignoring children
    Rule     4: HIT    : development.app1.service.url VALUE: 'jdbc:h2:/used'
    Rule     5: MISS   : development.app2.db.url, ignoring children

## 010 collisions
Sometimes you will create a situation where your path will resolve to multiple expressions of the same priority.  In this case, the engine won't know what to do.  Here is an example


    {
    	"development": {
    		"app1,app2.db.url" : "jdbc:h2:/first/path",
    		"app2,app3.db.url" : "jdbc:h2:/other/path",
    
    		"/app\\w+/.db.url" : "jdbc:h2:/regex/path",
    		"/ap\\w+/.db.url" : "jdbc:h2:/other_regex/path"
    	},
    	"development,qa": {
    		"example": "value_1"
    	},
    	"qa,production": {
    		"example": "value_2"
    	}
    }

The following examples will fail:

This is ambiguous, because it matches two enums

    $ go/confexpr-go lookup development.app2.db.url
    jdbc:h2:/other/path
This is ambiguous, because it matches two regexes

    $ go/confexpr-go lookup development.app4.db.url
    jdbc:h2:/regex/path
This is abiguous because it matches two enums at the top level

    $ go/confexpr-go lookup qa.example
    value_2
In the next section we'll see how to resolve the abiguity


The explain command can be used to gain insight into why the tool returned a certain value for an expression

This generates a collision at the enum

    $ go/confexpr-go explain development.app2.db.url
    The following rules were evaluated in this order, the first hit is returned
    Rule     1: HIT    : app2,app3.db.url VALUE: 'jdbc:h2:/other/path'
    Rule     2: COLLIDE: app1,app2.db.url VALUE: '<nil>'
    Rule     3: HIT    : /ap\w+/.db.url VALUE: 'jdbc:h2:/other_regex/path'
    Rule     4: COLLIDE: /app\w+/.db.url VALUE: '<nil>'
    Rule     5: HIT    : app2,app3.db.url VALUE: '<nil>'
    Rule     6: COLLIDE: app1,app2.db.url VALUE: '<nil>'
    Rule     7: HIT    : /ap\w+/.db.url VALUE: '<nil>'
    Rule     8: COLLIDE: /app\w+/.db.url VALUE: '<nil>'
    Rule     9: MISS   : qa,production, ignoring children
This generates a collision at the regex

    $ go/confexpr-go explain development.app4.db.url
    The following rules were evaluated in this order, the first hit is returned
    Rule     1: MISS   : app2,app3.db.url, ignoring children
    Rule     2: MISS   : app1,app2.db.url, ignoring children
    Rule     3: HIT    : /app\w+/.db.url VALUE: 'jdbc:h2:/regex/path'
    Rule     4: COLLIDE: /ap\w+/.db.url VALUE: '<nil>'
    Rule     5: MISS   : qa,production, ignoring children
    Rule     6: MISS   : app2,app3.db.url, ignoring children
    Rule     7: MISS   : app1,app2.db.url, ignoring children
    Rule     8: HIT    : /app\w+/.db.url VALUE: '<nil>'
    Rule     9: COLLIDE: /ap\w+/.db.url VALUE: '<nil>'
## 011 collision resolution
The proper way to resolve this abiguity for the paths in question is to provide a higher priority expression that is not ambiguous.  Here we are using direct matches to resolve the conflicts


    {
    	"development": {
    		"app1,app2.db.url" : "jdbc:h2:/first/path",
    		"app2,app3.db.url" : "jdbc:h2:/other/path",
    
    		"/app\\w+/.db.url" : "jdbc:h2:/regex/path",
    		"/ap\\w+/.db.url" : "jdbc:h2:/other_regex/path",
    
    		"app2.db.url" : "jdbc:h2:/direct/path",
    		"app4.db.url" : "jdbc:h2:/another_direct/path"
    	}
    }

This will produce the following output

    $ go/confexpr-go lookup development.app2.db.url
    'jdbc:h2:/direct/path'
     
    $ go/confexpr-go lookup development.app4.db.url
    'jdbc:h2:/another_direct/path'
     
In this example the regex patterns themselves are fairly poorly chosen, an you would probably want ot refactor the expressions

## 012 array access
We can also access array elements like any other element


    {
    	"development": {
    		"apps.dbs" : [
    			{
    				"user" : "app1_user",
    				"password" : "secret_1"
    			},
    			{
    				"user" : "app2_user",
    				"url" : "jdbc:h2:/other/path",
    				"password" : "secret_2"
    			},
    			{
    				"user" : "app3_user",
    				"password" : "secret_3"
    			}
    		],
    		"apps.dbs.*.url" : "jdbc:h2:/sample/path"
    
    	}
    }

This will produce the following output

    $ go/confexpr-go lookup development.apps.dbs.0.user
    'Error, state: State: Incomplete Path: [0 user] Evaluated_Path:[development apps dbs] Value:<nil> Variables:map[]'
     
    $ go/confexpr-go lookup development.apps.dbs.0.password
    'Error, state: State: Incomplete Path: [0 password] Evaluated_Path:[development apps dbs] Value:<nil> Variables:map[]'
     
You can see that wildcard substitution still works

    $ go/confexpr-go lookup development.apps.dbs.0.url
    'Error, state: State: Incomplete Path: [0 url] Evaluated_Path:[development apps dbs] Value:<nil> Variables:map[]'
     
And override rules are respected

    $ go/confexpr-go lookup development.apps.dbs.1.url
    'Error, state: State: Incomplete Path: [1 url] Evaluated_Path:[development apps dbs] Value:<nil> Variables:map[]'
     
# Expressions Reference
The following expressions are available in a label, with the high precedence matches towards the top

1. An exact match
1. An enum match
1. A regex match
1. A wildcard match
1. A deep wildcard match

## Exact Match
An exact match is simply a label.  It looks like this

    a.b.c

## An Enum Match
An enum match is delimited by commas.  It looks like this

	a1,a2.b.c

This will match `a1.b.c` and `a2.b.c`

## A Regex Match
A regex match is surrounded by slashes.  Using the `.` is NOT supported.  Instead, `\w` is your friend  

	/a\w*/.b.c

This will match `a.b.c`, `apple.b.c`, and `aardvark.b.c` 

## Wildcard Match
A wildcard match is specified like this

	*.b.c

This will match `a.b.c`, `b.b.c`, and any three element path that ends with `b.c`

## Deep Wildcard Match
A deep wildcard match is specified like this

	**.c

This will match anything ending with `.c`, such as `a.b.c`, `b.c`, or even `c`.  The expression is greedy, so it will also match `a.b.c.c`

# Learning More

There are collection of directories that include the actual source for this document,and they act as a test suite for the tool.  Each directory will contain:

* A config file, conf.jsonw
* A set of expected keys, in passing-keys.csv
* A human readable explanation of the test case, in input.md

You can play with these yourself, or add items to them.  Run the `tests.sh` command to regenerate this documentation.
