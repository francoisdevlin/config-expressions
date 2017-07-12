# 005 locality vs specificity
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
     
