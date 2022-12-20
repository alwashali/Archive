from flask_bcrypt import generate_password_hash
username=input("username:")
pw=input("password:")
pwd=generate_password_hash(pw)
print(pwd)


#add manually to sqlite explorer 
# INSERT INTO user VALUES (0,"admin","$2b$12$5/6OBq9DXlM5/GZkEVzK9eWjrA7u3vakVZ.DmmJm8xvkpUdh6NuHW");