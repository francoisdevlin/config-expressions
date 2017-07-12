This is a collection of tests that are run, to ensure that the rules are being interpolated properly.  Each directory is one test suite.  There are the following in each directory

* A config file, conf.jsonw
* A set of expected keys, in passing-keys.csv
* A human readable explaination of the test case, in README.md
## 001 vanilla json
This example shows that the tool can access configuration stored as vanilla json

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
    app1_user
     
    $ ./pattern-getter.rb development.app2.db.user
    app2_user
     
## 002 basic wildcard
This example show a basic wildcard match.  It is possible to use a wildcard to match both specified paths, such as app1.db, as well as dynamic paths, such as app9.db.url.  It also shows what happens in the case of a collision.  Notice that app2.db.url returns a different value, `jdbc:h2:/other/path`

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
    app1_user
     
    $ ./pattern-getter.rb development.app2.db.user
    app2_user
     
    $ ./pattern-getter.rb development.app1.db.url
    jdbc:h2:/sample/path
     
    $ ./pattern-getter.rb development.app2.db.url
    jdbc:h2:/other/path
     
    $ ./pattern-getter.rb development.app9.db.url
    jdbc:h2:/sample/path
     
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
    app1_user
     
    $ ./pattern-getter.rb development.app2.db.user
    app2_user
     
    $ ./pattern-getter.rb development.app1.db.url
    jdbc:h2:/sample/path
     
    $ ./pattern-getter.rb development.app2.db.url
    jdbc:h2:/sample/path
     
    $ ./pattern-getter.rb development.app9.db.user
    app9_user
     
    $ ./pattern-getter.rb development.app9.db.url
    jdbc:h2:/sample/path
     
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
    enum_app1_user
     
    $ ./pattern-getter.rb development.app2.db.user
    enum_app2_user
     
    $ ./pattern-getter.rb development.app1.db.url
    jdbc:h2:/sample/path
     
    $ ./pattern-getter.rb development.app2.db.url
    jdbc:h2:/sample/path
     
    $ ./pattern-getter.rb development.app9.db.user
    app9_user
     
    $ ./pattern-getter.rb development.app9.db.url
    jdbc:h2:/sample/path
     
## 005 locality vs specificity
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
    jdbc:h2:/sample/path
     
    $ ./pattern-getter.rb development.app2.db.url
    jdbc:h2:/other/path
     
## 006 hello

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

    $ ./pattern-getter.rb localhost.db.driver
    org.h2.Driver
     
    $ ./pattern-getter.rb localhost.sample.db.domain
    sample
     
    $ ./pattern-getter.rb dev02.sample.db.password
    dev002a_
     
    $ ./pattern-getter.rb dev02.sample.db.username
    sa
     
    $ ./pattern-getter.rb dev02.payment.db.username
    DDO_PAYMT_DBA_READ
     
    $ ./pattern-getter.rb localhost.sample.db.username
    sa
     
    $ ./pattern-getter.rb localhost.sample.db.password
    
     
    $ ./pattern-getter.rb qa01.sample.db.password
    deb09_Qa7
     
    $ ./pattern-getter.rb dev02.sample.db.url
    jdbc:oracle:thin:@dodcld.juniper.com:1521/ddebtomcatsvc
     
This is a collection of tests that are run, to ensure that the rules are being interpolated properly.  Each directory is one test suite.  There are the following in each directory

* A config file, conf.jsonw
* A set of expected keys, in passing-keys.csv
* A human readable explaination of the test case, in README.md
