# 001 vanilla json
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
     
