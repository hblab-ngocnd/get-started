function getDictionary(notCache){
    showLoading();
    let url = "./api/dictionary";
    let params = { not_cache:false, level:"n5" };
    if (notCache) {
        params.not_cache = true;
    }
    params.level = $('input[type=radio][name=oplevel]:checked').val()
    url = url + "?" + $.param(params);
    $.get(url)
        .done(function(data) {
            if(data.length > 0) {
                var content = "";
                $('#dict-table-body').html("");
                data.forEach(function(row,index) {
                    content = "<tr>";
                    content = content + "<td>" + (index + 1) + "</td>";
                    content = content + "<td>" + row.text + "</td>";
                    content = content + "<td>" + row.alphabet + "</td>";
                    content = content + "<td>" + row.mean_eng + "</td>";
                    content = content + "<td>" + row.mean_vn + "</td>";
                    content = content + "<td>" + row.detail + "</td>";
                    content = content + "</tr>";
                    $('#dict-table-body').append(content);
                });
                showPage();
            }
        })
        .fail(function() {
            showPage();
        })
}

$('input[type=radio][name=oplevel]').change(function() {
    getDictionary();
});

function showPage(){
    document.getElementById("loader").style.display = "none";
    document.getElementById("dict-table").style.display = "block";
}

function showLoading(){
    document.getElementById("loader").style.display = "block";
    document.getElementById("dict-table").style.display = "none";
}

//Call getDictionary on page load.
getDictionary();