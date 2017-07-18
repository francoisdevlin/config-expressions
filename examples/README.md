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

    $ ./config-expression lookup development.app1.db.user
    'app1_user'
     
    $ ./config-expression lookup development.app2.db.user
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

    $ ./config-expression lookup development.app1.db.user
    'app1_user'
     
    $ ./config-expression lookup development.app2.db.user
    'app2_user'
     
    $ ./config-expression lookup development.app1.db.url
    'jdbc:h2:/sample/path'
     
You can see that `app2.db.url` matches the more specific input

    $ ./config-expression lookup development.app2.db.url
    'jdbc:h2:/other/path'
     
We can also get a value for entries that aren't specified explicitly, such as `app9.db.url`

    $ ./config-expression lookup development.app9.db.url
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

    $ ./config-expression lookup development.app1.db.user
    'app1_user'
     
The url does not have any substitution, so no changes are made

    $ ./config-expression lookup development.app1.db.url
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

    $ ./config-expression lookup development.app1.db.user
    'enum_app1_user'
     
    $ ./config-expression lookup development.app2.db.user
    'enum_app2_user'
     
Notice that the wildcard pattern is matched after the enum is exhausted

    $ ./config-expression lookup development.app9.db.user
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

    $ ./config-expression lookup development.app1.db.user
    'regex_app1_user'
     
    $ ./config-expression lookup development.app2.db.user
    'regex_app2_user'
     
Notice the built in anchors to the regex, so this value falls through to the wildcard pattern

    $ ./config-expression lookup development.app10.db.user
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

    $ ./config-expression lookup development.app1.db.user
    'user'
     
    $ ./config-expression lookup development.app2.db.user
    'enum_app2_user'
     
    $ ./config-expression lookup development.app3.db.user
    'regex_app3_user'
     
    $ ./config-expression lookup development.app10.db.user
    'wildcard_app10_user'
     
Notice that the entire matched path is substituted

    $ ./config-expression lookup development.other_app.server.user
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

    $ ./config-expression lookup development.app1.db.url
    'jdbc:h2:/sample/path'
     
    $ ./config-expression lookup development.app2.db.url
    'jdbc:h2:/other/path'
     
Noting matches locally inside the `development` dictionary, so the engine falls back on the global value

    $ ./config-expression lookup development.app1.service.url
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

    $ ./config-expression explain development.app1.db.url
    The following rules were evaluated in this order, the first hit is returned
    Entering locality 'development'
    Rule     1: MISS   : development.app2
    Rule     2: HIT    : development.*.db.url VALUE: 'jdbc:h2:/sample/path'
    Entering locality ''
    Rule     3: HIT    : development.app1.db.url VALUE: 'jdbc:h2:/not/used'
    Rule     4: MISS   : development.app1.service
    Rule     5: MISS   : development.app2
    jdbc:h2:/sample/path
And here you can see there are no hits locally, and the global value is the highest priority hit

    $ ./config-expression explain development.app1.service.url jdbc:h2:/used
    The following rules were evaluated in this order, the first hit is returned
    Entering locality 'development'
    Rule     1: MISS   : development.app2
    Rule     2: MISS   : development.*.db
    Entering locality ''
    Rule     3: MISS   : development.app1.db
    Rule     4: HIT    : development.app1.service.url VALUE: 'jdbc:h2:/used'
    Rule     5: MISS   : development.app2
    jdbc:h2:/used
## 009 db connections
This is an example from real life.  The original file was approximately 360 lines of vanilla JSON.  This replacement version comes in at about 40 lines.  An order of magnitude improvement.  Not only is this a smaller file, but the real gains come when extending your system

* Adding new environments is a breeze, it will only require adding a top level url
* Adding a new database will usually require adding a user entry to the `*` locality

This drastically cuts down on the amount of busywork that is required for configuration


    {
    	"**.username":"sa",
    	"**.password":"",
    	"**.db.driver":"oracle.jdbc.OracleDriver",
    	"localhost,pipeline.db.driver" : "org.h2.Driver",
    	"*": {
    		"payment.db.username" : "DDO_PAYMT_DBA_READ",
    		"account.db.username" : "DDO_ACCT_DBA_READ",
    		"apply.db.username" : "DDO_APPLY_DBA_READ",
    		"communications.db.username" : "DDO_NTFCN_DBA_READ",
    		"idmap.db.username" : "DDO_IDMAP_DBA",
    		"kyc.db.username" : "DDO_KYC_DBA_READ",
    		"processorgateway.db.username" : "DDO_PRCSOR_DBA_READ",
    		"analyticsgateway.db.username" : "ODS_ANLYTC_GTWY_DBA"
    	},
    	"localhost":{
    		"*$domain.db.domain":"${domain}",
    		"analyticsgateway.db.domain" : "analytics-gateway",
    		"idmap.db.domain" : "platform-integration",
    		"processorgateway.db.domain" : "processor-gateway"
    	},
    	"ukdev" : {
    		"*.db.url" : "jdbc:oracle:thin:@//odu1-12-sl-uat.barcapint.com:3521/BCardNG.intranet.barcapint.com",
    		"payment.db.password" : "Barclays39#",
    		"account.db.password" : "Barclays09#",
    		"apply.db.password" : "Barclays12#",
    		"communications.db.password" : "Barclays33#",
    		"idmap.db.password" : "Barclays#123"
    	},
    
    	"dev02,dev03,qa03,cicluster.*.db.password" : "dev002a_",
    	"qa01,qa02,qa04.*.db.password" : "deb09_Qa7",
    
    	"pipeline.*.db.url" : "jdbc:h2:/local/domains/h2/payments;AUTO_SERVER=TRUE;MODE=Oracle" ,
    	"dev02.*.db.url" : "jdbc:oracle:thin:@dodcld.juniper.com:1521/ddebtomcatsvc",
    	"dev03.*.db.url" : "jdbc:oracle:thin:@dephcld.juniper.com:1521/deph3debtomcatsvc",
    	"qa01.*.db.url" : "jdbc:oracle:thin:@qodcld.juniper.com:1521/qdebtomcatsvc",
    	"qa02.*.db.url" : "jdbc:oracle:thin:@roracld.juniper.com:1521/qa02crdeb_svc",
    	"qa03.*.db.url" : "jdbc:oracle:thin:@qephcld.juniper.com:1521/qeph3debbatchsvc",
    	"qa04.*.db.url" : "jdbc:oracle:thin:@172.18.223.135:1521/DEVOPS01",
    	"cicluster.*.db.url" : "jdbc:oracle:thin:@dephcld.juniper.com:1521/deph3debtomcatsvc"
    }

This will produce the following output

    $ ./config-expression lookup localhost.db.driver
    'org.h2.Driver'
     
    $ ./config-expression lookup localhost.sample.db.domain
    'sample'
     
    $ ./config-expression lookup dev02.sample.db.password
    'dev002a_'
     
    $ ./config-expression lookup dev02.sample.db.username
    'sa'
     
    $ ./config-expression lookup dev02.payment.db.username
    'DDO_PAYMT_DBA_READ'
     
    $ ./config-expression lookup localhost.sample.db.username
    'sa'
     
    $ ./config-expression lookup localhost.sample.db.password
    ''
     
    $ ./config-expression lookup qa01.sample.db.password
    'deb09_Qa7'
     
    $ ./config-expression lookup dev02.sample.db.url
    'jdbc:oracle:thin:@dodcld.juniper.com:1521/ddebtomcatsvc'
     
### Further optimization

It is possible to optimize this event further with a small change to process.  For example, if we adopt a convention that every db user must be named after the application, we can replace the `*` dictionary with one entry



	"*.*$app_name.db.username" : "DDO_${app_name}_DBA_READ"



This will also remove the need to even maintain a list of usernames in _configuration_, as our pattern file will provide a _convention_ instead.

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

    $ ./config-expression lookup development.app2.db.url
    Ambiguous match for 'development.app2.db.url', the following expressions have the same priority:
    'development.app1,app2.db.url'
    'development.app2,app3.db.url'
    'development./app\w+/.db.url'
    'development./ap\w+/.db.url'
This is ambiguous, because it matches two regexes

    $ ./config-expression lookup development.app4.db.url
    Ambiguous match for 'development.app4.db.url', the following expressions have the same priority:
    'development./app\w+/.db.url'
    'development./ap\w+/.db.url'
This is abiguous because it matches two enums at the top level

    $ ./config-expression lookup qa.example
    Ambiguous match for 'qa.example', the following expressions have the same priority:
    'development,qa'
    'qa,production'
In the next section we'll see how to resolve the abiguity


The explain command can be used to gain insight into why the tool returned a certain value for an expression

This generates a collision at the enum

    $ ./config-expression explain development.app2.db.url
    The following rules were evaluated in this order, the first hit is returned
    Entering locality 'development'
    Rule     1: COLLIDE: development.app1,app2.db.url VALUE: ''
    Rule     2: COLLIDE: development.app2,app3.db.url VALUE: ''
    Rule     3: COLLIDE: development./app\w+/.db.url VALUE: ''
    Rule     4: COLLIDE: development./ap\w+/.db.url VALUE: ''
    Entering locality 'development,qa'
    Rule     5: MISS   : development,qa.example
    Entering locality ''
    Rule     6: MISS   : qa,production, ignoring children
    Could not find development.app2.db.url
This generates a collision at the regex

    $ ./config-expression explain development.app4.db.url
    The following rules were evaluated in this order, the first hit is returned
    Entering locality 'development'
    Rule     1: MISS   : development.app1,app2
    Rule     2: MISS   : development.app2,app3
    Rule     3: COLLIDE: development./app\w+/.db.url VALUE: ''
    Rule     4: COLLIDE: development./ap\w+/.db.url VALUE: ''
    Entering locality 'development,qa'
    Rule     5: MISS   : development,qa.example
    Entering locality ''
    Rule     6: MISS   : qa,production, ignoring children
    Could not find development.app4.db.url
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

    $ ./config-expression lookup development.app2.db.url
    'jdbc:h2:/direct/path'
     
    $ ./config-expression lookup development.app4.db.url
    'jdbc:h2:/another_direct/path'
     
In this example the regex patterns themselves are fairly poorly chosen, an you would probably want ot refactor the expressions

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
