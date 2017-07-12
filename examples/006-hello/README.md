# 006 hello

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
     