const httpRoot = "http://localhost:294";

function convert() {
    let input = document.getElementById("input");
    let output = document.getElementById("output");
    let inputModel = {
        ConvertType: "DecToHex",
        InputData: input.value
    }
    fetch(httpRoot + "/convert", {
        method: "POST",
        body: JSON.stringify(inputModel),
        headers: {"Content-Type": "application/json;charset=UTF-8",},
    }).then(re => {
        if (re.ok) return re.json();
    }).then(re => {
        output.value = re.Data.OutputData;
        updateMessage(re)
    })
}

function loadRecord() {
    let selectRecord = document.getElementById("selectRecord");
    fetch(httpRoot + "/loadRecord?RecordId=" + selectRecord.value, {
        method: "GET",
        headers: {"Content-Type": "application/json;charset=UTF-8",},
    }).then(re => {
        if (re.ok) return re.json();
    }).then(re => {
        updateMessage(re)
        if (re.Success) {
            let selectGroup = document.getElementById("selectGroup");
            let selectType = document.getElementById("selectType");
            let input = document.getElementById("input");
            let output = document.getElementById("output");
            selectGroup.value = re.Data.GroupId;
            selectType.value = re.Data.ConvertType;
            input.value = re.Data.InputData;
            output.value = re.Data.OutputData;
        }
        output.value = re.Data.OutputData;
    })
}

function saveRecord() {
    let recordName = prompt("Input a name", "");
    if (recordName == null || recordName == "") {
        alert("The name cannot be empty")
        return;
    }
    let selectGroup = document.getElementById("selectGroup");
    let selectType = document.getElementById("selectType");
    let input = document.getElementById("input");
    let output = document.getElementById("output");
    let record = {
        Id: 0,
        Name: recordName,
        ConvertType: selectType.value,
        InputData: input.value,
        OutputData: output.value,
        GroupId: Number(selectGroup.value),
    }
    fetch(httpRoot + "/saveRecord", {
        method: "POST",
        body: JSON.stringify(record),
        headers: {"Content-Type": "text/html; charset=utf-8",},
    }).then(re => {
        if (re.ok) return re.text();
    }).then(re => {
        let divHeader = document.getElementById("divHeader");
        divHeader.innerHTML = re;
        let messageElement = document.getElementById("message");
        messageElement.innerText = "Record was saved!";
    })
}

function deleteRecord() {
    let selectRecord = document.getElementById("selectRecord");
    fetch(httpRoot + "/deleteRecord?RecordId=" + selectRecord.value, {
        method: "DELETE",
        headers: {"Content-Type": "text/html; charset=utf-8",},
    }).then(re => {
        if (re.ok) return re.text();
    }).then(re => {
        let divHeader = document.getElementById("divHeader");
        divHeader.innerHTML = re;
        let messageElement = document.getElementById("message");
        messageElement.innerText = "Record was deleted!";
    })
}

function createGroup() {

}

function deleteGroup() {

}

function updateMessage(re) {
    if (re.Success) {
        let messageElement = document.getElementById("message");
        messageElement.innerText = re.Message;
    } else {
        alert(re.Message)
    }
}