include("dbconfig.php");
include_once("include/errstr.php");
include_once("include/ucall.php");
function stncall($call) {
  $ucall = strtoupper($call);
  $p = strpos($ucall,"/P");
  if($p > 2) {
    $ucall = substr($ucall,0,$p);
  }
  return $ucall;
}

function act($user,$summit,$con) {
  $r = false;
  $sql = "SELECT * FROM `activator_log` WHERE `activatedby` = '".$user."' AND `wotaid` = ".$summit;
  $res = mysql_query($sql,$con);
  if($res && (mysql_num_rows($res) > 0)) {
    $r = true;
  }
  return $r;
}

function act_yr($user,$year,$summit,$con) {
  $r = false;
  $sql = "SELECT * FROM `activator_log` WHERE `activatedby` = '".$user."' AND `year` = ".$year." AND `wotaid` = ".$summit;
  $res = mysql_query($sql,$con);
  if($res && (mysql_num_rows($res) > 0)) {
    $r = true;
  }
  return $r;
}

function logcontact($call,$summitid,$activator,$callused,$date,$s2s,$con) {
  global $valid_cnt;
  if($call!=="") {
    $tz = new DateTimeZone('GMT');
    $dt = new DateTime($date, $tz);
    $year = substr($date,0,4);

// replace $wotaid = "LDW-".substr("000".$summitid,-3,3);

$fellid = (int)$summitid;
   if ($fellid <= 214) {$prefix = "LDW-";}
  else {$prefix = "LDO-";
   $fellid = (int)$fellid-214;}
   $ldid = $prefix . str_pad($fellid,3,"0",STR_PAD_LEFT);
   $wotaid = $ldid;

    echo "<p>Contact with $call from $wotaid on ".$dt->format('D j M Y')." ";
    // check it is not already in log
    $sql = "SELECT * FROM `activator_log` WHERE `activatedby` = '".$activator."' AND `wotaid` = '".$summitid."' AND `stncall` = '".$call."' AND `date` = '".$date."'";
    $result = mysql_query($sql,$con);
    if($result && (mysql_num_rows($result) > 0)) {
      echo "already in log.</p>";
    } else {
      // check confirmation
      $sql = "SELECT * FROM `chaser_log` WHERE `ucall` = '".$call."' AND `wotaid` = '".$summitid."' AND `stncall` = '".$callused."' AND `date` = '".$date."'";
      $result = mysql_query($sql,$con);
      if($result && (@mysql_num_rows($result) > 0)) {
        $cfm = "true";
        // update record
        mysql_query("UPDATE `chaser_log` SET `confirmed` = true WHERE `ucall` = '".$call."' AND `wotaid` = '".$summitid."' AND `stncall` = '".$callused."' AND `date` = '".$date."'");
        echo "Confirmed ";
      } else {
        $cfm = "false";
      }
      // add contact
      $sql = "INSERT into `activator_log` (`activatedby`, `callused`, `wotaid`, `date`, `year`, `stncall`, `ucall`, `s2s`, `confirmed`) VALUES ('".$activator."','".$callused."',".$summitid.",'".$date."',".$year.",'".$call."','".ucall($call)."',".$s2s.",".$cfm.")";
      $result = mysql_query($sql,$con);
      echo "Logged</p>";
      $valid_cnt++;
    }
  }
}
$self = $_SERVER['PHP_SELF']."?page=".$_REQUEST['page'];
// initialization
global $gCms;
$feu =& $gCms->modules['FrontEndUsers']['object'];
$logged_by = $feu->LoggedInName();
$callused = ucall($logged_by)."/P";
global $errcount;
$errcount = 0;
global $valid_cnt;
$valid_cnt = 0;
if($_SERVER['REQUEST_METHOD'] == "POST") {
  // load the form fields
  $date = isset($_POST['date1']) ? mysql_real_escape_string($_POST['date1']) : "";
  $tz = new DateTimeZone('GMT');
  $dt = new DateTime($date, $tz);
  $year = substr($date,0,4);
  $summitid = mysql_real_escape_string($_POST['summit']);

// replace $wotaid = "LDW-".substr("000".$summitid,-3,3);

 $fellid = (int)$summitid;
   if ($fellid <= 214) {$prefix = "LDW-";}
  else {$prefix = "LDO-";
   $fellid = (int)$fellid-214;}
   $ldid = $prefix . str_pad($fellid,3,"0",STR_PAD_LEFT);
   $wotaid = $ldid;

  $activator = strtoupper(trim(mysql_real_escape_string($_POST['logged_by'])));
  $callused = stncall(strtoupper(trim(mysql_real_escape_string($_POST['call_used']))));
  $call_1 = strtoupper(trim(mysql_real_escape_string($_POST['contact_1'])));
  $s2s_1 = isset($_POST['s2s_1']) ? "true" : "false";
  $call_2 = strtoupper(trim(mysql_real_escape_string($_POST['contact_2'])));
  $s2s_2 = isset($_POST['s2s_2']) ? "true" : "false";
  $call_3 = strtoupper(trim(mysql_real_escape_string($_POST['contact_3'])));
  $s2s_3 = isset($_POST['s2s_3']) ? "true" : "false";
  $call_4 = strtoupper(trim(mysql_real_escape_string($_POST['contact_4'])));
  $s2s_4 = isset($_POST['s2s_4']) ? "true" : "false";
  $cnt = isset($_POST['cnt']) ? trim(mysql_real_escape_string($_POST['cnt'])) : "0";
  $valid_cnt = isset($_POST['vcnt']) ? trim(mysql_real_escape_string($_POST['vcnt'])) : "0";
  if($call_1=="") {
    $message = errstr("No contacts specified");
  }
  switch($year) {
    case '2009': break;
    case '2010': break;
    case '2011': break;
    case '2012': break;
    case '2013': break;
    case '2014': break;
    case '2015': break;
    case '2016': break;
    case '2017': break;
    case '2018': break;
    case '2019': break;
    default: $message = errstr("Cannot log contacts for year ".$year);
  }
}
if(($_SERVER['REQUEST_METHOD'] == "POST")&&($errcount==0)) {
  $act = act($activator,$summitid,$con);
  $act_yr = act_yr($activator,$year,$summitid,$con);
  logcontact($call_1,$summitid,$activator,$callused,$date,$s2s_1,$con);
  logcontact($call_2,$summitid,$activator,$callused,$date,$s2s_2,$con);
  logcontact($call_3,$summitid,$activator,$callused,$date,$s2s_3,$con);
  logcontact($call_4,$summitid,$activator,$callused,$date,$s2s_4,$con);
  if($valid_cnt >= 1) {
    if(!$act_yr) {
      echo "<p>Activator points for $year table awarded: 1</p>";
    }
    if(!$act) {
        echo "<p>Activator points for all-time table awarded: 1</p>";
    }
    // set last activator details in summit table
    mysql_query("UPDATE `summits` SET `last_act_by` = '".$activator."', `last_act_date` = '".$date."' WHERE `wotaid` = '".$summitid."'");
  }
  $summitid = $summitid + 0;
?>
<br>
<h3>Log more contacts</h3>
<form action="<?=$self?>" method="POST">
    <input type="hidden" name="logged_by" value="<?=$logged_by?>">
    <input type="hidden" name="call_used" value="<?=$callused?>">
    <input type="hidden" name="date1" value="<?=$date?>">
    <input type="hidden" name="summit" value="<?=$summitid?>">
    <table border="0" cellpadding="2">
        <tr>
            <td width="25%" style="padding-top:4px;padding-bottom:4px;">Date of activation:&nbsp;</td>
            <td><?=$dt->format('j M Y')?></td>
        </tr>
        <tr>
            <td width="25%" style="padding-top:4px;padding-bottom:4px;">Summit:</td>
            <td width="75%"><?=$wotaid?></td>
        </tr>
<?
} else {
  require_once('classes/tc_calendar.php');
  //instantiate class and set properties
  $myCalendar = new tc_calendar("date1", true);
  $myCalendar->setIcon("images/iconCalendar.gif");
  $myCalendar->setDate(date(d),date(m),date(Y));
  // set contact count
  $cnt = 0;
  // show error message, if any
  if($errcount > 0) echo $message;
  // display the form
?>
<script language="javascript" src="calendar.js"></script>
<form action="<?=$self?>" method="POST">
    <input type="hidden" name="logged_by" value="<?=$logged_by?>">
    <table border="0" cellpadding="2">
        <tr>
            <td>Date of activation:&nbsp;</td>
            <td><? $myCalendar->writeScript(); ?></td>
        </tr>
        <tr>
            <td>Callsign used:&nbsp;</td>
            <td><input type="text" size="8" name="call_used" value="<?=$callused?>"></td>
        </tr>
        <tr>
            <td width="25%">Summit:</td>
            <td width="75%">
<select name="summit">
<?
$result = mysql_query("SELECT `wotaid`,`sotaid`,`name` FROM `summits` WHERE `wotaid` > 0 ORDER BY `name`",$con);
if (mysql_num_rows($result) ) {
  while($row = mysql_fetch_array($result, MYSQL_ASSOC)) {
    $val = $row['wotaid'] + 0;

  //  $wotaid = "LDW-".substr($row['wotaid'],-3,3);

$fellid1 = substr($row['wotaid'],-3,3);
  $fellid = (int)$fellid1;
  if ($fellid <= 214) {$prefix = "LDW-";}
  else {$prefix = "LDO-";
  $fellid = (int)$fellid-214;}
  $ldid = $prefix . str_pad($fellid,3,"0",STR_PAD_LEFT);
 $wotaid = $ldid;
 $sotaid= "";
  if ($row['sotaid'] != '') {
   $sotano = (int)$row['sotaid'];
   $sotaref = "G/LD-" . str_pad($sotano,3,"0",STR_PAD_LEFT);
   $sotaid = "[". $sotaref . "]";
  }
    echo "<option value=\"".$val."\">".$row['name']." (".$wotaid.") ".$sotaid." </option>";
  }
}
?>
</select>
            </td>
<?
}
?>
        </tr>
        <tr>
            <td>Contact <? $cnt++; echo $cnt;?>:</td>
            <td><input type="text" size="8" name="contact_1">&nbsp;&nbsp;<input type="checkbox" name="s2s_1" value="Y"> Summit-to-summit? </td>
        </tr>
        <tr>
            <td>Contact <? $cnt++; echo $cnt;?>:</td>
            <td><input type="text" size="8" name="contact_2">&nbsp;&nbsp;<input type="checkbox" name="s2s_2" value="Y"> Summit-to-summit? </td>
        </tr>
        <tr>
            <td>Contact <? $cnt++; echo $cnt;?>:</td>
            <td><input type="text" size="8" name="contact_3">&nbsp;&nbsp;<input type="checkbox" name="s2s_3" value="Y"> Summit-to-summit? </td>
        </tr>
        <tr>
            <td>Contact <? $cnt++; echo $cnt;?>:</td>
            <td><input type="text" size="8" name="contact_4">&nbsp;&nbsp;<input type="checkbox" name="s2s_4" value="Y"> Summit-to-summit? </td>
        </tr>
    </table>
    <p><input type="submit" value="Submit">&nbsp;&nbsp;<button type="button" onclick="javascript:window.location='mm_home.html'">Quit</button></p>
  <input type="hidden" name="cnt" value="<?=$cnt?>">
  <input type="hidden" name="vcnt" value="<?=$valid_cnt?>">
</form>
<?