The goal of this document is to show both the _why_ as much as the _how_ to use this pattern matching library.  It starts by walking you through an example of a static, vanilla JSON file, and adds more an more features to add more an more power.  At the end is a reference section on each pattern

# Learning by Example
Let's start with a very basi example, using vanilla JSON
## 001 vanilla json
Vanillia JSON is a configuration format we are all used to using.  There is a lot of pre-existing configuration written as vanilla JSON, so reaing it directly made sense as a starting point for the tool.  It is where you will start when trying to abstract out patterns from pre-existing files


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

    $ ./pattern-getter.rb development.app1.db.user
    'app1_user'
     
    $ ./pattern-getter.rb development.app2.db.user
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

    $ ./pattern-getter.rb development.app1.db.user
    'app1_user'
     
    $ ./pattern-getter.rb development.app2.db.user
    'app2_user'
     
    $ ./pattern-getter.rb development.app1.db.url
    'jdbc:h2:/sample/path'
     

You can see that `app2.db.url` matches the more specific input


    $ ./pattern-getter.rb development.app2.db.url
    'jdbc:h2:/other/path'
     

We can also get a value for entries that aren't specified explicitly, such as `app9.db.url`


    $ ./pattern-getter.rb development.app9.db.url
    'jdbc:h2:/sample/path'
     
## 003 basic substitution
This example shows to to use basic variable subsitution with a wildcard.  It matches keys that are static, such as app1.db, as well as keys that are dynamic, such as app9.


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

    $ ./pattern-getter.rb development.app1.db.user
    'app1_user'
     
    $ ./pattern-getter.rb development.app2.db.user
    'app2_user'
     
    $ ./pattern-getter.rb development.app1.db.url
    'jdbc:h2:/sample/path'
     
    $ ./pattern-getter.rb development.app2.db.url
    'jdbc:h2:/sample/path'
     
    $ ./pattern-getter.rb development.app9.db.user
    'app9_user'
     
    $ ./pattern-getter.rb development.app9.db.url
    'jdbc:h2:/sample/path'
     
## 004 basic enum
This example shows to to use an enum, with variable substition.  You'll notice the following:

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

    $ ./pattern-getter.rb development.app1.db.user
    'enum_app1_user'
     
    $ ./pattern-getter.rb development.app2.db.user
    'enum_app2_user'
     
    $ ./pattern-getter.rb development.app1.db.url
    'jdbc:h2:/sample/path'
     
    $ ./pattern-getter.rb development.app2.db.url
    'jdbc:h2:/sample/path'
     
    $ ./pattern-getter.rb development.app9.db.user
    'app9_user'
     
    $ ./pattern-getter.rb development.app9.db.url
    'jdbc:h2:/sample/path'
     
## 005 basic regex
The enum matcher is a valuable way to use a union type.  However, you must explicitly include every match you want to include.  Sometimes it is more useful to match a more general pattern, such as a regular expression.  This example shows a regex matcher at work.  The regex is specified in the pattern `/app\\w/$app_name.db.user`. Notice the following

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

    $ ./pattern-getter.rb development.app1.db.user
    'regex_app1_user'
     
    $ ./pattern-getter.rb development.app2.db.user
    'regex_app2_user'
     

Notice the built in anchors to the regex, so this value falls through to the wildcard pattern


    $ ./pattern-getter.rb development.app10.db.user
    'app10_user'
     
## 006 complete precedence
This example shows the complete precedence of patterns of the same locality.  The order is

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

    $ ./pattern-getter.rb development.app1.db.user
    'user'
     
    $ ./pattern-getter.rb development.app2.db.user
    'enum_app2_user'
     
    $ ./pattern-getter.rb development.app3.db.user
    'regex_app3_user'
     
    $ ./pattern-getter.rb development.app10.db.user
    'wildcard_app10_user'
     

Notice that the entire matched path is substituted


    $ ./pattern-getter.rb development.other_app.server.user
    'deep_wildcard_other_app.server_user'
     
## 007 locality vs specificity
Locality is established by creating a new dictionary object

This example shows local instructions winning oevr global ones.  Observe development.app1.db.url.  Even thought the global rule is more specific, the specification of a local wildcard rule overrides the global one.  This is because the most local rules win a conflict. 

You can see that more specific rules still win if they have the same locality, in development.app2.db.url


    {
    	"development.app1.db.url" : "jdbc:h2:/not/used",
    
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

    $ ./pattern-getter.rb development.app1.db.url
    'jdbc:h2:/sample/path'
     
    $ ./pattern-getter.rb development.app2.db.url
    'jdbc:h2:/other/path'
     
## 008 hello


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


Hello


    $ ./pattern-getter.rb localhost.db.driver
    'org.h2.Driver'
     
    $ ./pattern-getter.rb localhost.sample.db.domain
    'sample'
     
    $ ./pattern-getter.rb dev02.sample.db.password
    'dev002a_'
     
    $ ./pattern-getter.rb dev02.sample.db.username
    'sa'
     
    $ ./pattern-getter.rb dev02.payment.db.username
    'DDO_PAYMT_DBA_READ'
     
    $ ./pattern-getter.rb localhost.sample.db.username
    'sa'
     
    $ ./pattern-getter.rb localhost.sample.db.password
    ''
     
    $ ./pattern-getter.rb qa01.sample.db.password
    'deb09_Qa7'
     
    $ ./pattern-getter.rb dev02.sample.db.url
    'jdbc:oracle:thin:@dodcld.juniper.com:1521/ddebtomcatsvc'
     
# Patterns Reference
The following patterns are available in a label, with the high precedence matches towards the top

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

This will match anything ending with `.c`, such as `a.b.c`, `b.c`, or even `c`.  The pattern is greedy, so it will also match `a.b.c.c`

# Learning More

There are collection of directories that include the actual source for this document,and they act as a test suite for the tool.  Each directory will contain:

* A config file, conf.jsonw
* A set of expected keys, in passing-keys.csv
* A human readable explaination of the test case, in input.md

You can play with these yourself, or add items to them.  Run the `tests.sh` command to regenerate this documentation.
