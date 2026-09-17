"use strict";
const query=document.getElementById("query");
query.addEventListener("input",()=>{
 const list=document.getElementById("results");list.replaceChildren();
 const term=query.value.trim().toLocaleLowerCase("ru");
 const hits=term?window.BOOK_SEARCH.filter(x=>x.text.toLocaleLowerCase("ru").includes(term)):[];
 for(const hit of hits){const li=document.createElement("li");const a=document.createElement("a");a.href=hit.url;a.textContent=hit.title;li.append(a);list.append(li);}
 document.getElementById("search-status").textContent=term?"Найдено глав: "+hits.length:"";
});
