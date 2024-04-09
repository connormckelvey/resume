### Work Experience {pStyle="SectionHeading" class="SectionHeading"}

<%
    work.map(item => template("work-item.md", item))
        .join("\n")
%>
