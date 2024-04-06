<% 
    (function() {
        const months = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"]; 

        function formatDate(str) {
            const date = new Date(str)        
            const year = date.getFullYear();
            const monthName = months[date.getMonth()];
            return`${monthName} ${year}`;
        }

        return `${formatDate(dates.start)} - ${dates.end ? formatDate(dates.end) : "Present" }`
    })() 
%>