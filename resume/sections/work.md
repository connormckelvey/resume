### Work Experience {pStyle="SectionHeading"}

<%
    work.map(item => template("work-item.md", item))
        .join("\n")
%>
