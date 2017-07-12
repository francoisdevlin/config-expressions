# 004 basic enum
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
     
