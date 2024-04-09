### Projects {pStyle="SectionHeading" class="SectionHeading"}

<%
    projects
        .map((project) => {
            const link = project.links && project.links[0]
                ? project.links[0]
                : null
            const name = link
                ? `[**${project.name}**](${link.href})`
                : `**${project.name}**`
            return template("../util/definition-list.md", { 
                term: name, 
                definition: project.description
            })
        })
        .join("\n\n")
%>
