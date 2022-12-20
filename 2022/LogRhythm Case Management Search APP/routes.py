from app import create_app,db,login_manager,bcrypt
from models import User
import requests,time,datetime
import pandas as pd 
from forms import login_form
from dotenv import load_dotenv
import os

from flask import (
    render_template,
    redirect,
    flash,
    request,
    url_for,
    session,
    send_from_directory
)

from datetime import timedelta

from flask_bcrypt import Bcrypt,generate_password_hash, check_password_hash

from flask_login import (
    
    login_user,
    logout_user,
    login_required,
)

from app import create_app,db,login_manager,bcrypt
from models import User
from forms import login_form



df = pd.DataFrame
token=''
Cases_API='https://IPAddress:8501/lr-case-api/cases/'



app = create_app()


@login_manager.user_loader
def load_user(user_id):
    return User.query.get(int(user_id))


@app.before_request
def session_handler():
    session.permanent = True
    app.permanent_session_lifetime = timedelta(minutes=480)


@app.route("/login/", methods=("GET", "POST"), strict_slashes=False)
def login():
    form = login_form()

    if form.validate_on_submit():
        try:
            user = User.query.filter_by(username=form.username.data).first()
            if check_password_hash(user.pwd, form.pwd.data):
                login_user(user)
                return redirect(url_for('index'))
            else:
                flash("Invalid Username or password!", "danger")
        except Exception as e:
            flash(e, "danger")

    return render_template("auth.html",
        form=form,
        text="Login",
        title="Login",
        btn_action="Login"
        )


@app.route('/downloads/<filename>')
@login_required
def download(filename):

    if filename != "":

        urlmetric = Cases_API
        headermetric = {'Authorization': 'Bearer ' + token,
                        'Content-Type': 'application/json', }

        for i in df.index:
            metricresponse = None
            metricresponse = requests.get(urlmetric + str(df['number'][i]) + "/metrics/", headers=headermetric,verify=False)

            if metricresponse.status_code == 200:
                df.loc[i, 'FirstAlarmTriggered'] = metricresponse.json()['earliestEvidence']['originalDate']
                time.sleep(0.100)

        
        # Change timezone and format 
        df['FirstAlarmTriggered'] = pd.to_datetime(df.FirstAlarmTriggered, format='%Y-%m-%d %H:%M:%S')
        df['FirstAlarmTriggered']= df['FirstAlarmTriggered'].dt.tz_convert('Asia/Riyadh') 
        df['FirstAlarmTriggered']= df['FirstAlarmTriggered'].dt.strftime('%Y-%m-%d %H:%M')
        

        
        casesdata = df[['id', 'number','name','summary','priority','severity','dateCreated','FirstAlarmTriggered','status.name','dateUpdated', 'dateClosed','owner.name','collaborators']]
        
        casesdata.to_excel('./downloads/'+filename)

        full_path = os.path.join(app.root_path, "downloads/")

        return send_from_directory(full_path, filename)
    else:
        return render_template('home.html',error="error: No file to download. Please make sure to search first!")

def getdata1(clientname, createdafter):
    global token
    userinputcreatedafter = (datetime.datetime.now() - datetime.timedelta(days=int(createdafter))).isoformat('T') + 'Z'
    load_dotenv('.env')
    if clientname == 'Client1':
        token=os.environ.get('Entity1')
    elif clientname == 'Client2':
        token=os.environ.get('Entity2')
    elif clientname == 'Client3':
        token=os.environ.get('Entity2')
    else:
        print("Error : no token found")
        return "notoken"

    urlcases = Cases_API
    headercases = { 'Authorization': 'Bearer '+token,'Content-Type': 'application/json', 'createdAfter': userinputcreatedafter,'count':'9000'}
    response = requests.get(urlcases, headers=headercases, verify=False)
    # Access to 
    if response.status_code == 401:
        
        return render_template('home.html', error="Error:Not Authorised, Check Token permission")
    else:

        global df
        df = pd.json_normalize(response.json(), max_level=1)
        df.loc[df['priority'] == 1, 'severity'] = 'Critical'
        df.loc[df['priority'] == 2, 'severity'] = 'High'
        df.loc[df['priority'] == 3, 'severity'] = 'Medium'
        df.loc[df['priority'] == 4, 'severity'] = 'Low'
        df.loc[df['priority'] == 5, 'severity'] = 'Informational'

        #Change to Datetime type 
        df['dateCreated']= pd.to_datetime(df.dateCreated, format='%Y-%m-%d %H:%M')  
        df['dateUpdated']= pd.to_datetime(df.dateUpdated, format='%Y-%m-%d %H:%M')  

        # Change timezone
        df['dateUpdated']= df['dateUpdated'].dt.tz_convert('Asia/Riyadh') 
        df['dateCreated']= df['dateCreated'].dt.tz_convert('Asia/Riyadh') 

        # Change datetime format 
        df['dateUpdated']= df['dateCreated'].dt.strftime('%Y-%m-%d %H:%M')
        df['dateCreated']= df['dateCreated'].dt.strftime('%Y-%m-%d %H:%M')



        tempdf = df[['number', 'name', 'summary', 'dateCreated','dateUpdated', 'owner.name', 'priority','status.name','severity']]

        #reverse order
        tempdf = tempdf[::-1]  
    
        userfilename = 'LRCases-'+ clientname + datetime.datetime.now().strftime("%d-%m-%Y-%H-%S")
        downloadfilename=userfilename+'.xlsx'

        temp = tempdf.to_dict('records')

        columnnames=tempdf.columns.values
    
        return temp,columnnames,downloadfilename


@login_required
@app.route('/getdata', methods=['POST','GET'])
def getdata():
    if request.method == "POST":
        clienttoken = request.form['clientname']
        searchdays= request.form['searchdays']
        if len(clienttoken) == 0 or len(searchdays) == 0:
            return render_template('home.html',error="Error:PLease Select Client and Time")
        else:
            result = getdata1(clienttoken, searchdays)
            return render_template('home.html',records=result[0], colnames=result[1],downloadfilename=result[2])
    else:
        return render_template('home.html')



@app.route('/', methods=['GET','POST'])
@login_required
def home():
    return render_template('home.html')


@app.route('/index', methods=['GET','POST'])
@login_required
def index():
    return render_template('home.html')

@app.route('/register', methods=['GET','POST'])
@login_required
def register():
    return render_template('404.html')

@app.route("/logout")
@login_required
def logout():
    logout_user()
    return redirect(url_for('login'))

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=5000)
