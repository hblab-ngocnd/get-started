function getDictionary(notCache){
    showLoading();
    let url = "./api/dictionary";
    let params = { not_cache:false, level:"n5", start:0,page_size: 20};
    if (notCache) {
        params.not_cache = true;
    }
    params.page_size = parseInt(document.getElementById("paging").getAttribute("data-page-size"),10)
    params.start = parseInt(document.getElementById("paging").getAttribute("data-start"),10)
    params.level = $('input[type=radio][name=oplevel]:checked').val()
    url = url + "?" + $.param(params);
    $.get(url)
        .done(function(data) {
            if(data.length > 0) {
                let content = "";
                if(params.start == 0){
                    $('#dict-table-body').html("");
                }
                data.forEach(function(row,index) {
                    content = "<tr>";
                    content = content + "<td>" + (index + params.start + 1) + "</td>";
                    content = content + "<td>" + row.text + "</td>";
                    content = content + "<td>" + row.alphabet + "</td>";
                    content = content + "<td>" + row.mean_eng + "</td>";
                    content = content + "<td>" + row.mean_vn + "</td>";
                    content = content + "<td><div class='content-detail' id='content-detail-"+index+"'>" + row.detail + "</div><button class='btn btn-default btn-sm btn-detail' data-toggle='modal' data-target='.bd-example-modal-lg' data-detail='content-detail-"+index+"' onclick='showDetail(this)'><span class='glyphicon glyphicon-info-sign'></span> 詳細</button></td>";
                    content = content + "</tr>";
                    $('#dict-table-body').append(content);
                });
                let pageSize = parseInt(document.getElementById("paging").getAttribute("data-page-size"),10)
                let prevStart = parseInt(document.getElementById("paging").getAttribute("data-start"),10)
                document.getElementById("paging").setAttribute("data-start",(prevStart+pageSize).toString(10))
                document.getElementById("paging").setAttribute("data-lock","false");
            }
            showPage();
        })
        .fail(function() {
            showPage();
        })
}

window.onscroll = function(ev) {
    if ((window.innerHeight + window.scrollY) >= document.body.offsetHeight) {
        if(document.getElementById("paging").getAttribute("data-lock") == "false"){
            document.getElementById("paging").setAttribute("data-offset",document.body.offsetHeight.toString());
            getDictionary();
        }
    }
};

$('input[type=radio][name=oplevel]').change(function() {
    reload();
});

function reload(){
    document.getElementById("paging").setAttribute("data-lock","false");
    document.getElementById("paging").setAttribute("data-start","0");
    document.getElementById("paging").setAttribute("data-offset","0");
    getDictionary();
}

function showDetail(el){
    let a = "<div class='modal-detail'>" + document.getElementById($(el).attr("data-detail")).innerHTML + "</div>";
   document.getElementsByClassName("modal-content-detail")[0].innerHTML = a;
}

function showPage(){
    document.getElementById("loader").style.display = "none";
    document.getElementById("dict-table").style.display = "block";
    document.body.scrollTop = parseInt(document.getElementById("paging").getAttribute("data-offset"),10);
}

function showLoading(){
    document.getElementById("loader").style.display = "block";
    document.getElementById("dict-table").style.display = "none";
    document.getElementById("paging").setAttribute("data-lock","true");
}

//Call getDictionary on page load.
getDictionary();