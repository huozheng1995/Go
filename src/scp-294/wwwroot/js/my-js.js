const httpRoot = "http://localhost:294";

document.addEventListener("DOMContentLoaded", () => {
    changeGroup();
});

function convert() {
    let input = document.getElementById("input");
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
        setOutput(re.Data.OutputData);
        updateMessage(re)
    })
}

function clearText() {
    let input = document.getElementById("input");
    let output = document.getElementById("output");
    input.value = null;
    output.value = null;
    updateMessageValue(null);
}

function loadRecord() {
    let selectRecord = document.getElementById("selectRecord");
    fetch(httpRoot + "/loadRecord?RecordId=" + selectRecord.value, {
        method: "GET",
        headers: {"Content-Type": "application/json;charset=UTF-8",},
    }).then(re => {
        if (re.ok) return re.json();
    }).then(re => {
        if (re.Success) {
            let selectGroup = document.getElementById("selectGroup");
            let selectType = document.getElementById("selectType");
            let input = document.getElementById("input");
            selectGroup.value = re.Data.GroupId;
            selectType.value = re.Data.ConvertType;
            input.value = re.Data.InputData;
            setOutput(re.Data.OutputData);
        }
        updateMessage(re)
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
        reloadGroup(groupId);
        updateMessageValue("Record was saved!")
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
        reloadGroup(groupId);
        updateMessageValue("Record was deleted!")
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
        updateMessageValue("Group was added!")
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
        updateMessageValue("Group was deleted!")
    })
}

function reloadGroup(groupId) {
    let selectGroup = document.getElementById("selectGroup");
    for (let option of selectGroup.options) {
        if (option.value == groupId) {
            option.selected = true;
            changeGroup();
            break;
        }
    }
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
    if (!selected) {
        selectRecord.selectedIndex = -1;
    }
}

function setOutput(data) {
    let output = document.getElementById("output");
    output.value = data;
    if (output.scrollHeight > 200) {
        output.style.height = '200px';
        output.style.height = output.scrollHeight + 64 + 'px';
    }
}

function updateMessage(re) {
    let element = document.getElementById("message");
    element.innerText = re.Message;
    element.style.color = re.Success ? "green" : "red";
}

function updateMessageValue(value) {
    let element = document.getElementById("message");
    element.innerText = value;
    if (value != null) {
        element.style.color = "green";
    }
}