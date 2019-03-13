List:
    GET /api/systems
    Request: <empty>
    Response:
        {
            "id":1,
            "name":"QA-FUNC-11",
            "description":null,
            "uuid":"62cadcb6-94eb-48ab-aca0-2d28b9492b36",
            "status":"down",
            "connection_status":"connection_ok",
            "uptime":null,
            "replication_level":2,
            "control_address":"localhost",
            "control_port":10016,
            "nfs_address":"192.168.0.1",
            "nfs_ip_range":24,
            "data_address":"10.0.0.1",
            "data_ip_range":16,
            "data_vlan":null,
            "external_use_dhcp":null,
            "external_address":null,
            "external_ip_range":null,
            "external_gateway":null,
            "external_network":null,
            "created_at":"2015-11-15T09:37:06.000Z",
            "updated_at":"2015-11-24T15:24:26.000Z"
        }

Login:
    POST /api/sessions
    Request:
        {"user":{"login":"admin","password":"changeme"}}
    Response:
        {"info":"Logged in","user":{"id":1,"login":"admin","admin":true,"first_name":"Super","surname":"Admin","email":"admin@example.com","created_at":"2015-11-15T09:37:06.000Z","updated_at":"2015-11-24T15:24:44.551Z"}}


Setup:
	Request: 
        {"skip_tests":false}
	Response:
        {"id":1,"name":"QA-FUNC-11","description":null,"uuid":"62cadcb6-94eb-48ab-aca0-2d28b9492b36","status":"down","connection_status":"connection_ok","uptime":null,"replication_level":2,"control_address":"localhost","control_port":10016,"nfs_address":"192.168.0.1","nfs_ip_range":24,"data_address":"10.0.0.1","data_ip_range":16,"data_vlan":null,"external_use_dhcp":null,"external_address":null,"external_ip_range":null,"external_gateway":null,"external_network":null,"created_at":"2015-11-15T09:37:06.000Z","updated_at":"2015-11-24T15:24:26.000Z"}

Start:
	Request:
        {"create_defaults":false}
    Response:
        {"id":1,"name":"QA-FUNC-11","description":null,"uuid":"62cadcb6-94eb-48ab-aca0-2d28b9492b36","status":"in_service","connection_status":"connection_ok","uptime":"2015-11-24T15:38:23.000Z","replication_level":2,"control_address":"localhost","control_port":10016,"nfs_address":"192.168.0.1","nfs_ip_range":24,"data_address":"10.0.0.1","data_ip_range":16,"data_vlan":null,"external_use_dhcp":null,"external_address":null,"external_ip_range":null,"external_gateway":null,"external_network":null,"created_at":"2015-11-15T09:37:06.000Z","updated_at":"2015-11-24T15:38:23.000Z"}

