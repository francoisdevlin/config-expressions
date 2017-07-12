# 002 basic wildcard
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
     
