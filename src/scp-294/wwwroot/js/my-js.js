const httpRoot = "http://localhost:294";

function convert() {
    console.log("start")
    let input = document.getElementById("input");
    let output = document.getElementById("output");
    let inputModel = {
        ConvertType: "DecToHex",
        InputData: input.value
    }
    fetch(httpRoot + "/convert", {
        method: "POST",
        body: JSON.stringify(inputModel),
        headers: {
            "Content-Type": "application/json;charset=UTF-8",
        },
    }).then((re) => {
        if (re.ok) {
            return re.json();
        }
    }).then((re) => {
        output.value = re.Data.OutputData;
        updateMessage(re.Message)
    })
}

function loadRecord() {
    console.log("abc")
    alert("aaa")
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
    console.dir(record)
    fetch(httpRoot + "/saveRecord", {
        method: "POST",
        body: JSON.stringify(record),
        headers: {
            "Content-Type": "application/json;charset=UTF-8",
        },
    }).then((re) => {
        if (re.ok) {
            return re.json();
        }
    }).then((re) => {
        updateMessage(re.Message)
    })

}


function updateMessage(message) {
    let messageElement = document.getElementById("message");
    messageElement.innerText = message;
}