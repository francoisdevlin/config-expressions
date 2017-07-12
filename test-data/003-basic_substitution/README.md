# 003 basic substitution
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
     
