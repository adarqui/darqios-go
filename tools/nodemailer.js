var nodemailer = require("nodemailer");

if(!process.env['SUBJECT'] || !process.env['BODY']) process.exit(-1);

var smtpTransport = nodemailer.createTransport("SMTP",{
    service: "Gmail",
    auth: {
        user: "critical@host.com",
        pass: "password"
    }
});

var mailOptions = {
    from: "critical@host.com",
    to: "alerts@host.com",
    cc : ["someone@host.com"],
    subject: process.env['SUBJECT'],
    text: process.env['BODY'],
}

smtpTransport.sendMail(mailOptions, function(error, response){
    if(error){
        console.log(error);
        process.exit(-1);
    }else{
        console.log("Message sent: " + response.message);
        process.exit(0);
    }
});
