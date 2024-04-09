### Education {pStyle="SectionHeading" class="SectionHeading"}

<%
    education.map(item => template("education-item.md", item))
        .join("\n")
%>
