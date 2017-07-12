Locality is established by creating a new dictionary object

This example shows local instructions winning oevr global ones.  Observe development.app1.db.url.  Even thought the global rule is more specific, the specification of a local wildcard rule overrides the global one.  This is because the most local rules win a conflict. 

You can see that more specific rules still win if they have the same locality, in development.app2.db.url
