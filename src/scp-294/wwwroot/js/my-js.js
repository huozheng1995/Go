const httpRoot = "http://localhost:294";

document.addEventListener("DOMContentLoaded", () => {
    changeGroup();
});

function convert() {
    let input = document.getElementById("input");
    let output = document.getElementById("output");
    let selectType = document.getElementById("selectType");
    let inputModel = {
        ConvertType: selectType.value,
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

function clearText() {
    let input = document.getElementById("input");
    let output = document.getElementById("output");
    input.value = null;
    output.value = null;
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
            let selectType = document.getElementById("selectType");
            let input = document.getElementById("input");
            let output = document.getElementById("output");
            selectType.value = re.Data.ConvertType;
            input.value = re.Data.InputData;
            output.value = re.Data.OutputData;
        }
    })
}

function addRecord() {
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
    let groupId = document.getElementById("selectGroup").value;
    fetch(httpRoot + "/addRecord", {
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
        let selectGroup = document.getElementById("selectGroup");
        for (let option of selectGroup.options) {
            if (option.value == groupId) {
                option.selected = true;
                changeGroup();
                break;
            }
        }
    })
}

function deleteRecord() {
    if (!confirm("Delete it?")) {
       return;
    }
    let selectRecord = document.getElementById("selectRecord");
    let groupId = document.getElementById("selectGroup").value;
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
        let selectGroup = document.getElementById("selectGroup");
        for (let option of selectGroup.options) {
            if (option.value == groupId) {
                option.selected = true;
                changeGroup();
                break;
            }
        }
    })
}

function createGroup() {
    let groupName = prompt("Input a name", "");
    if (groupName == null || groupName == "") {
        alert("The name cannot be empty")
        return;
    }
    let group = {
        Id: 0,
        Name: groupName,
    }
    fetch(httpRoot + "/addGroup", {
        method: "POST",
        body: JSON.stringify(group),
        headers: {"Content-Type": "text/html; charset=utf-8",},
    }).then(re => {
        if (re.ok) return re.text();
    }).then(re => {
        let divHeader = document.getElementById("divHeader");
        divHeader.innerHTML = re;
        let messageElement = document.getElementById("message");
        messageElement.innerText = "Group was added!";
    })
}

function deleteGroup() {
    if (!confirm("Delete it?")) {
        return;
    }
    let selectGroup = document.getElementById("selectGroup");
    fetch(httpRoot + "/deleteGroup?GroupId=" + selectGroup.value, {
        method: "DELETE",
        headers: {"Content-Type": "text/html; charset=utf-8",},
    }).then(re => {
        if (re.ok) return re.text();
    }).then(re => {
        let divHeader = document.getElementById("divHeader");
        divHeader.innerHTML = re;
        let messageElement = document.getElementById("message");
        messageElement.innerText = "Group was deleted!";
    })
}

function changeGroup() {
    let selectGroup = document.getElementById("selectGroup");
    let selectRecord = document.getElementById("selectRecord");
    let selected = false;
    for (let option of selectRecord.options) {
        option.selected = false;
        if (option.getAttribute("groupId") == selectGroup.value) {
            option.style.display = null;
            if (!selected) {
                option.selected = true;
                selected = true;
            }
        } else {
            option.style.display = "none";
        }
    }
}

function updateMessage(re) {
    if (re.Success) {
        let messageElement = document.getElementById("message");
        messageElement.innerText = re.Message;
    } else {
        alert(re.Message)
    }
}